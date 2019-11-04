[![Go Report Card](https://goreportcard.com/badge/edgexfoundry-holding/edgex-cli)](https://goreportcard.com/report/edgexfoundry-holding/edgex-cli)

# edgex-cli

## Introduction

A command line interface to interact with EdgeX microservices. Replaces the need to manually construct complex CURL commands and/or maintain developer scripts.

## Installation

In order to run this tool, you will need a running EdgeX instance and Go 1.12 or higher installed on your machine.

* Clone the git repo:

```
$ git clone https://github.com/edgexfoundry-holding/edgex-cli
```

* Install the CLI:

```
$ make install
```

You can now use the CLI by entering `edgex-cli` anywhere on your machine. 

### Available Commands

#### Readings

```
$ edgex reading list    // All readings
$ edgex reading list <device> // All readings for device, default limit is 50
$ edgex reading list <device> -l <limit>  // All readings for device with limit
```

#### Notifications

```
$ edgex notification add </path/to/file.toml> // Add notifications from toml file
$ edgex notification rm <slug> // Remove a notification by slug
$ edgex notification rm --age 1561939200 // Remove a notification by age (timestamp)
```

#### Subscriptions

```
$ edgex subscription add </path/to/file.toml> // Add subscriptions from toml file
$ edgex subscription rm A // Remove a subscription by slug
$ edgex subscription add /path/to/file.toml // Add subscriptions from toml file.
$ edgex subscription rm A // Remove a subscription by slug.
$ edgex subscription rm --id SOMEID // Remove a subscription by id.
```

#### Intervals
```
$ edgex interval list // Lists all intervals provided their numbers is within the configured limit.
$ edgex interval list A // List interval with name A.
$ edgex interval list --id SOMEID // List interval with id SOMEID.
$ edgex interval add /path/to/file.toml // Add intervals from toml file.
$ edgex interval update /path/to/file.toml // Update intervals from toml file.
$ edgex interval rm A // Remove a interval by id or name A.
```

### Developers

To try out your changes, you can either build the binary or calling `go run`.

* Build and run:

```
$ make build
$ ./edgex-cli
```

* Use `go run`:

```
$ go run main.go [COMMAND]
```

* Running tests:

```
$ make test
```

This will generate a coverage.out file in the root directory of the repo which you can use to see test coverage of code by running

```
$ go tool cover -html=coverage.out
```