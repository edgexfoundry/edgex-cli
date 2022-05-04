# EdgeX CLI Snap
[![edgex-cli](https://snapcraft.io/edgex-cli/badge.svg)][edgex-cli]

This directory contains the snap packaging of the EdgeX CLI.

The snap is built automatically and published on the Snap Store as [edgex-cli].

For usage instructions, please refer to EdgeX CLI section in [Getting Started using Snaps][docs].

## Build from source
Execute the following command from the top-level directory of this repo:
```
snapcraft
```

This will create a snap package file with `.snap` extension. It can be installed locally by setting the `--dangerous` flag:
```bash
sudo snap install --dangerous <snap-file>
```

The [snapcraft overview](https://snapcraft.io/docs/snapcraft-overview) provides additional details.

[edgex-cli]: https://snapcraft.io/edgex-cli
[docs]: https://docs.edgexfoundry.org/2.2/getting-started/Ch-GettingStartedSnapUsers/#edgex-cli
