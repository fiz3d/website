# [![Fiz3D](http://fiz3d.org/static/media/readme_small.png)](https://fiz3d.org) [Website](https://fiz3d.org/) [![Build Status](https://circleci.com/gh/fiz3d/website.svg?&style=shield)](https://circleci.com/gh/fiz3d/website)

This repository holds everything that runs on the [fiz3d.org](https://fiz3d.org) website.

## Dependencies

Below is a list of software used (some of which you will need installed):

- [Go](https://golang.org) is used to power the back-end.
- [Bootstrap](https://getbootstrap.com) and [jQuery](https://jquery.com) power the front-end.
- [Browserify](https://browserify.org) and [browserify-css](https://www.npmjs.com/package/browserify-css) are used to package resources, with [watchify](https://www.npmjs.com/package/watchify) for live development.
- [Rego](https://sourcegraph.com/github.com/sqs/rego) automatically rebuilds the Go source code as you make changes.

## Development

- Install the dependencies listed above (or `sudo make deps`).
- `cd path/to/this/repo`
- `make`

## Deployment

The app itself runs with a full repository clone, and is fully responsible for building itself. In practice, this is pretty fast.

Once any change is merged into `master`, the live site will pull those changes and trigger a rebuild / relaunch itself (disable this behavior with `-update=false`).

#### Simple

The simplest way to deploy is to use the `ubuntu-install.sh` script. This installs all relevant dependencies (git, Go, etc) and runs `make provision` for you:

```
wget -O - https://raw.github.com/fiz3d/website/master/ubuntu-install.sh | bash
```

#### Manual

`make provision` sets up a machine to run the site (and `make unprovision` reverts the changes to your system), so you can `service fiz3d-org start`.

#### Debugging

Check `/var/log/upstart/fiz3d-org.log` for details about what is going on.

If you notice the service doesn't start or receive a cryptic error like:

```
Failed to start fiz3d-org.service: Unit fiz3d-org.service failed to load: No such file or directory.
```

Your system probably uses systemd (e.g. Ubuntu 15.04 does) and not upstart (we currently only have an upstart script available).

## Issues

Please file all issues on the primary [fiz repository](https://github.com/fiz3d/fiz): [create a new issue](https://github.com/fiz3d/fiz/issues/new).

## License

#### General

- All source code in this repository (scripts, Go source code, etc) is licensed under the 3-clause BSD license (see the `LICENSE` file).
- All media in this repository (image, CSS, HTML, HTML templates, etc) are licensed under the CC BY 4.0 (Creative Commons) license, which can be [viewed online](https://creativecommons.org/licenses/by/4.0/).

#### Exceptions

The above does not apply for the contents of directories that explicitly contain their own licenses, listed below for your convenience:

- Fonts in `static/media/fonts` (see `OFL.txt`).
