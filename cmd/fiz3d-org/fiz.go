package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	addr       = flag.String("http", ":8080", "HTTP address to serve on")
	tlsAddr    = flag.String("https", "", "HTTPS address to serve on")
	pemFile    = flag.String("pem-file", "", "HTTPS / SSL certificate .pem file")
	keyFile    = flag.String("key-file", "", "HTTPS / SSL certificate .key file")
	update     = flag.Bool("update", true, "update via Git and shutdown server after pull")
	updateRate = flag.String("update-rate", "60s", "rate at which to check for updates via Git")
	dev        = flag.Bool("dev", false, "reload all templates on each request")

	staticDirPrefix = "/static/" // Path prefix to strip to get static dir root.
	staticDir       = "static/"  // Directory to serve for static files.

	errorTemplate = "error.tmpl" // Template to use for errors.
	templateExt   = ".tmpl"      // Filepath extension of template files.
	templateIncl  = ".incl"      // Filepath extension of template includes.
	templateDir   = "templates/" // Relative directory that templates reside in.

	src = &GitUpdater{
		Dir: ".",
	}
	lastUpdate  = time.Now()
	tmpls       *template.Template
	updateRateT time.Duration
)

func reloadTemplates() error {
	// Recursively match all filepaths with the template extension.
	var matches []string
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if e := filepath.Ext(path); e == templateExt || e == templateIncl {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Load the templates with proper names.
	tmpls = template.New("")
	for _, m := range matches {
		name := strings.TrimPrefix(m, templateDir)
		data, err := ioutil.ReadFile(m)
		if err != nil {
			return err
		}
		if _, err := tmpls.New(name).Parse(string(data)); err != nil {
			return err
		}
	}
	return nil
}

type TemplateData struct {
	Request *http.Request

	Error string // Error in serving the request, if any.
	Stack string // Stack trace, if any.
}

func handler(w http.ResponseWriter, r *http.Request) error {
	if *dev {
		if err := reloadTemplates(); err != nil {
			return err
		}
	}

	// Determine template name.
	tmplName := r.URL.Path[1:]
	if stat, err := os.Stat(filepath.Join(templateDir, tmplName)); err != nil {
		return err
	} else {
		if stat.IsDir() {
			tmplName = path.Join(tmplName, "index")
		}
	}
	tmplName = tmplName + ".tmpl"

	// Find the requested page.
	root := tmpls.Lookup(tmplName)
	if root == nil {
		// Couldn't find the page, render the error page.
		w.WriteHeader(http.StatusNotFound)
		return fmt.Errorf("%v - %v", http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	// Execute the requested template page. Buffer the output so that we can
	// handle template errors without corrupting the response.
	buf := bytes.NewBuffer(nil)
	err := root.Execute(buf, &TemplateData{
		Request: r,
	})
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	return err
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v\n", r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

var stackTrace = &struct {
	sync.Mutex
	buf []byte
}{
	buf: make([]byte, 8192),
}

func stack() string {
	stackTrace.Lock()
	defer stackTrace.Unlock()

read:
	n := runtime.Stack(stackTrace.buf, false)
	if n > len(stackTrace.buf) {
		stackTrace.buf = make([]byte, len(stackTrace.buf)*2)
		goto read
	}
	return string(stackTrace.buf)
}

func errorHandler(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	render := func(w http.ResponseWriter, r *http.Request, err, stack string) {
		log.Printf("[error] %v %v - %v\n", r.Method, r.URL, err)

		// Render the error template.
		data := &TemplateData{
			Request: r,
			Error:   err,
			Stack:   stack,
		}
		if err := tmpls.ExecuteTemplate(w, errorTemplate, data); err != nil {
			panic(err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from panics to render them as an error.
		defer func() {
			if err := recover(); err != nil {
				render(w, r, fmt.Sprintf("%v", err), stack())
			}
		}()

		// Execute the handler for the request, render the error if any.
		err := h(w, r)
		if err == nil {
			log.Printf("%v %v\n", r.Method, r.URL)
			return
		}
		render(w, r, err.Error(), "")
	})
}

func checkForUpdates() {
	// If not enough time has elapsed, skip the update.
	if time.Since(lastUpdate) < updateRateT {
		return
	}
	lastUpdate = time.Now()

	// Check for updates.
	updated, err := src.Update()
	if err != nil {
		log.Println("Update error:", err)
		return
	}
	if updated {
		log.Println("Updated source code. Exiting server..")
		os.Exit(0)
	}
	log.Println("checked for updates (none available)")
}

func init() {
	if err := reloadTemplates(); err != nil {
		log.Fatal(err)
	}

	if *update {
		go func() {
			for {
				time.Sleep(updateRateT)
				checkForUpdates()
			}
		}()
	}
}

func main() {
	flag.Parse()

	// Parse update rate flag.
	var err error
	updateRateT, err = time.ParseDuration(*updateRate)
	if err != nil {
		log.Fatal(err)
	}

	// Static file hosting
	http.Handle("/static/", logHandler(http.StripPrefix(staticDirPrefix, http.FileServer(http.Dir(staticDir)))))

	// App handler.
	http.Handle("/", errorHandler(handler))

	// Start HTTPS server:
	if *tlsAddr != "" {
		go func() {
			if *pemFile == "" {
				log.Fatal("expected -pem-file flag")
			}
			if *keyFile == "" {
				log.Fatal("expected -key-file flag")
			}
			log.Println("Serving on", *tlsAddr)
			log.Fatal(http.ListenAndServeTLS(*tlsAddr, *pemFile, *keyFile, nil))
		}()
	}

	// Start HTTP server:
	log.Println("Serving on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
