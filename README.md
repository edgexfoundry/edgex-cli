[![Go Report Card](https://goreportcard.com/badge/edgexfoundry/edgex-cli)](https://goreportcard.com/report/edgexfoundry/edgex-cli)

# edgex-cli

## Introduction

A command line interface to interact with EdgeX microservices. Replaces the need to manually construct complex CURL commands and/or maintain developer scripts.

```
 ______     _              __   __            _____  _       _____        
|  ____|   | |             \ \ / /           / ____|| |     |_   _|     
| |__    __| |  __ _   ___  \ V /   ______  | |     | |       | |        
|  __|  / _` | / _` | / _ \  > <   |______| | |     | |       | |        
| |____| (_| || (_| ||  __/ / . \           | |____ | |____  _| |_       
|______|\__,_| \__, | \___|/_/ \_\           \_____||______||_____| 
		__/ |                                                             
	       |___/                                                              

EdgeX CLI version: 0.0.1
https://www.edgexfoundry.org/

Usage:
  edgex-cli [command]

Available Commands:
  addressable   Addressable command command
  command       `Command` command
  db            Purges entire EdgeX Database. [USE WITH CAUTION]
  device        Device command
  deviceservice Device service command
  event         Event command
  help          Help about any command
  interval      Interval command
  notification  Notification command
  profile       Device profile command.
  reading       Reading command
  status        Checks the current status of each microservice.
  subscription  Subscription command
  version       Version command

Flags:
      --config-file string   configuration file
  -h, --help                 help for edgex-cli
      --no-pager             Do not pipe output into a pager.
  -u, --url                  Print URL(s) used by the entered command.
  -v, --verbose              Print entire HTTP response.

Use "edgex-cli [command] --help" for more information about a command.
```

## Installation

In order to run this tool, you will need an accessible **running EdgeX instance** running locally, accessible via localhost.
 and Go 1.12 or higher installed on your machine.

* Clone the git repo:

```
$ git clone https://github.com/edgexfoundry/edgex-cli
```

* Change directory:

```
$ cd edgex-cli
```

* Install the CLI:

```
$ make install
```
Install also makes a copy of the default configuration and copies it to $HOME/.edgex-cli/configuration.toml.
If your EdgeX instance is not running on localhost, minimally you will need to replace localhost with the correct IP address.
You can now use the CLI by entering `edgex-cli` anywhere on your machine provided your $GOBIN is on your $PATH.


#### Interactive Mode

Some commands leverage interactive-mode which opens an editor and allows you to provide information that would 
normally be difficult with just command line arguments. For example, creating an Event requires a lot of information,
also Events contain zero or more readings. Using interactive mode, you can easily create an Event with many readings and 
customize each reading. You can choose the editor that is used by setting the environment variable `EDITOR`. The default
editor is `Vi` for `Unix` operating systems(MacOS, Linux, etc), `Notepad` for Windows OS. The default editor is used if 
no `EDITOR` is specified. Some examples of editors:

- vi
- vim
- nano
- emacs
- notepad
- vscode
- atom


### Developers

Follow the installation instructions to create a configuration file and edit as necessary.
To try out your changes, you can either build the binary and execute or by calling `go run`.

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

#### Code Organization
All CLI go code lives under the "cmd" directory, who sub-directories map one is to one to the toplevel commands, namely
addressable, command, db, device, ... version.
The children of the above map to the sub-commands supported by each of these.
Use the -f flag to provide file input, while n is reserved for names of objects.
There are command specific and global flags such as verbose respectively. 
Read the Cobra and Viper documentation for additional help.


#### Sample templates
Sample template files are in the "samples" directory for device profile (createDP.json and yaml),
device (createDevice.toml), intervals (createInterval.toml and json) and for updateInterval.
Edit these should data structures change. Add more should you want to ease defining other objects.
Occasionally samples are multiple formats, this was to support for legacy formats. Going forward we may 
only support json format.

## Supported commands and sub-commands

```
    addressable
      list        A list of all addressable
```   
```
    command
      list        A list of device supported commands
```   
```
    db
      purge - Purges entire EdgeX Database. [USE WITH CAUTION]
```

```
    device
      add         Add devices
      adminstate  Device admin state 
      list        A list of all devices
      rm          Removes device by name or ID
      update      Update a device
```

```
    deviceservice
      add         Add a device service
      list        Lists existing devices services
      rm          Removes device service by name or ID
      update      Update device service 
```
```
    event
      add         Create an event
      count       Returns the count of core-data events
      list        A list of Events
      rm          Removes event by its id or removes all events generated by given device
      scrub       Remove all (pushed) events and their associated readings [USE WITH CAUTION]
 ```         
```
    interval
      add         Add interval
      list        A list of all intervals
      rm          Removes interval by name or id
      update      Update interval

``` 
```
    notification
      add         Add notification
      list        A list of all notifications
      rm          Removes notification by slug or age
```  
```
    profile
      add         Add profiles
      list        Returns a list of device profiles
      rm          Remove profile by name or ID
```  
```
    reading
      count       Returns the count of core-data readings
      list        A list of readings across devices or pertaining to a specified device
```  
```
    status        This command pings each edgex microservice and prints their status.
                  This command is not stable yet.
``` 
```
    subscription
      add         Add subscription
      list        A list of all subscriptions
      rm          Removes subscription by --slug or id.
```  
```
    version       Version command
```
