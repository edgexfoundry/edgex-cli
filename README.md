# edgex-cli

## Introduction

A command line interface to interact with EdgeX microservices. You can use this CLI tool instead of complicated CURL command or developer scripts.

## Installation

In order to run this tool, you will need a running EdgeX instance and Go 1.11 or higher installed on your machine.

* Clone the git repo in your go path:

```
$ cd ~/go/src/github.com/edgexfoundry/
$ git clone https://github.com/edgexfoundry-holding/edgex-cli
```

*A bug described in [issue #26](https://github.com/edgexfoundry-holding/edgex-cli/issues/26) is solved through the use of the replace directive. This fix, however, can be erased when using go tidy and some other go commands. make install and make build work.*

* Install the CLI:

```
$ make install
```

~~You can now use the CLI by entering `edgex` anywhere on your machine.~~ *BUG: `make install` currently installs the binary as `edgex-cli` globally, because of the directory structure.*

### Developers

To test your changes, you can either build the binary or calling `go run`.

* Build and run:

```
$ make build
$ ./edgex
```

* Use `go run`:

```
$ go run main.go [COMMAND]
```
