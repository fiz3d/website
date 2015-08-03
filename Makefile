dev:
	rego github.com/fiz3d/website/cmd/fiz3d-org $(FLAGS)

provision: unprovision
	cp upstart.conf /etc/init/fiz3d-org.conf
	chown root:root /etc/init/fiz3d-org.conf
	chmod 644 /etc/init/fiz3d-org.conf
	ln -s /lib/init/upstart-job /etc/init.d/fiz3d-org

unprovision:
	rm -f /etc/init/fiz3d-org.conf
	rm -f /etc/init.d/fiz3d-org
