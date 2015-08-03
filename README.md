# The fiz3d.org website

This repository holds everything that runs on the [fiz3d.org](https://fiz3d.org) website.

## Dependencies

Below is a list of software used (some of which you will need installed):

- [Go](https://golang.org) is used to power the back-end.
- [Bootstrap](https://getbootstrap.com) and [jQuery](https://jquery.com) power the front-end.
- [Browserify](https://browserify.org) is used to package resources.
- [Docker](https://www.docker.com) is used for quick and easy deployment.
- [Travis CI](https://travis-ci.org) lets us auto-deploy any changes merged into this repository directly to the [live site](https://fiz3d.org).
- [Rego](https://sourcegraph.com/github.com/sqs/rego) automatically rebuilds the Go source code as you make changes.

## Development

- Install the dependencies listed above.
- `cd path/to/this/repo`
- `make`

## Deployment

The app itself runs with a full repository clone, and is fully responsible for building itself. In practice, this is pretty fast.

Once any change is merged into `master`, the live site will pull those changes and trigger a rebuild / relaunch itself (disable this behavior with `-update=false`).

`make provision` sets up a machine to run the site (and `make unprovision` reverts the changes to your system), so you can `service fiz3d-org start`.

If you notice the service doesn't start or receive a cryptic error like:

```
Failed to start fiz3d-org.service: Unit fiz3d-org.service failed to load: No such file or directory.
```

Your system probably uses systemd (e.g. Ubuntu 15.04 does) and we currently only offer an upstart script.

## License

- All source code in this repository (scripts, Go source code, etc) is licensed under the 3-clause BSD license (see the `LICENSE` file).
- All media in this repository (image, CSS, HTML, HTML templates, etc) are licensed under the CC BY 4.0 (Creative Commons) license, which can be [viewed online](https://creativecommons.org/licenses/by/4.0/).

Except for the contents of directories that explicitly contain their own licenses (listed below for your convenience):

- Fonts in `static/media/fonts` (see `OFL.txt`).
