#!/bin/sh -e

# configuration.toml is required for the client to run
if [ ! -f "$SNAP_USER_DATA/.edgex-cli/configuration.toml" ]; then
    mkdir -p "$SNAP_USER_DATA/.edgex-cli"
    cp "$SNAP/res/configuration.toml" "$SNAP_USER_DATA/.edgex-cli/configuration.toml"
fi

exec "$@"
