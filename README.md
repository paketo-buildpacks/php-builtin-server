# Paketo PHP Built-in Web Server Cloud Native Buildpack
The PHP Built-in Web Server Buildpack sets up and runs PHP's built-in web
server to support PHP applications

The buildpack is published for consumption at `gcr.io/paketo-buildpacks/php-builtin-server` and
`paketo-buildpacks/php-builtin-server`.

## Behavior
This buildpack is the default web-server in the PHP buildpack.
If the `BP_PHP_SERVER` environment variable is set to `php-server` at
build-time this buildpack will participate. It will also participate if the
environment variable isn't set at all.

The buildpack will do the following:
* At run time:
  - Contribute a start command for the PHP built-in webserver

This buildpack `requires` `php` at launch time, and will also optionally
`require` `composer-packages` at launch-time if the application contains a
`composer.json` file or the `$COMPOSER` environment variable is set.

## Configuration

### `BP_PHP_WEB_DIR`
The web directory or document root can be configured via the `BP_PHP_WEB_DIR`
environment variable. Set the environment variables at build time either
directly  or through a [`project.toml`
file](https://github.com/buildpacks/spec/blob/main/extensions/project-descriptor.md).

#### `pack build` flag
```shell
pack build my-app --env BP_PHP_WEB_DIR="htdocs"
```

#### In a [`project.toml`](https://github.com/buildpacks/spec/blob/main/extensions/project-descriptor.md)
```toml
[build]
  [[build.env]]
    name = 'BP_PHP_WEB_DIR'
    value = 'htdocs'
```

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh --version <version-number>
```

This will create a `buildpackage.cnb` file under the `build` directory which you
can use to build your app as follows:
`pack build <app-name> -p <path-to-app> -b build/buildpackage.cnb -b <other-buildpacks..>`

To run the unit and integration tests for this buildpack:
```
$ ./scripts/unit.sh && ./scripts/integration.sh
```
