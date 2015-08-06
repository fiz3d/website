.PHONY: dev clean deps provision unprovision

dev: clean
	cd static/js/ && watchify -t browserify-css site.js -o site.min.js &
	rego github.com/fiz3d/website/cmd/fiz3d-org -dev -update=false $(FLAGS)

clean:
	rm -f ./static/js/site.min.js

deps:
	npm install -g browserify
	npm install -g watchify

provision: unprovision
	cp upstart.conf /etc/init/fiz3d-org.conf
	chown root:root /etc/init/fiz3d-org.conf
	chmod 644 /etc/init/fiz3d-org.conf
	ln -s /lib/init/upstart-job /etc/init.d/fiz3d-org
	service fiz3d-org start

unprovision:
	rm -f /etc/init/fiz3d-org.conf
	rm -f /etc/init.d/fiz3d-org
