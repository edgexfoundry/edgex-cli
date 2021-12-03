# EdgeX CLI Snap
[![snap store badge](https://raw.githubusercontent.com/snapcore/snap-store-badges/master/EN/%5BEN%5D-snap-store-black-uneditable.png)](https://snapcraft.io/edgex-cli)

This folder contains snap packaging for the EdgeX-CLI Snap

The snap currently supports both `amd64` and `arm64` platforms.

## Installation

### Installing snapd
The snap can be installed on any system that supports snaps. You can see how to install 
snaps on your system [here](https://snapcraft.io/docs/installing-snapd/6735).

However for full security confinement, the snap should be installed on an 
Ubuntu 18.04 LTS or later (Desktop or Server), or a system running Ubuntu Core 18 or later.

### Installing EdgeX CLI as a snap
The snap is published in the snap store at https://snapcraft.io/edgex-cli.
You can see the current revisions available for your machine's architecture by running the command:

```bash
$ snap info edgex-cli
```

The latest stable version of the snap can be installed using:

```bash
$ sudo snap install edgex-cli
```

The latest development version of the snap can be installed using:

```bash
$ sudo snap install edgex-cli --edge
```

A specific version of the snap can be installed by setting the channel, for instance for 2.1 (Jakarta):

```bash
$ sudo snap install edgex-cli --channel=2.1
```

For the older CLI version compatible with 1.x of EdgeX, use `--channel=1.0`



**Note** - the snap has only been tested on Ubuntu Core, Desktop, and Server.
