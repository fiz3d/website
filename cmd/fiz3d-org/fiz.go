package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addr       = flag.String("http", ":8080", "HTTP address to serve on")
	tlsAddr    = flag.String("https", "", "HTTPS address to serve on")
	pemFile    = flag.String("pem-file", "", "HTTPS / SSL certificate .pem file")
	keyFile    = flag.String("key-file", "", "HTTPS / SSL certificate .key file")
	update     = flag.Bool("update", true, "update via Git and shutdown server after pull")
	updateRate = flag.String("update-rate", "60s", "rate at which to check for updates via Git")

	src = &GitUpdater{
		Dir: ".",
	}
	tmpls       = template.Must(template.ParseGlob("templates/*"))
	lastUpdate  = time.Now()
	updateRateT time.Duration
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v\n", r.Method, r.URL)
	fmt.Fprintf(w, "Hello world!")
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// App handler.
	http.HandleFunc("/", handler)

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
