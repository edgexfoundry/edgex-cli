#!/bin/sh -e

# configuration.toml is required for the client to run
if [ ! -f "$HOME/.edgex-cli/configuration.toml" ]; then
    mkdir -p "$HOME/.edgex-cli"
    cp "$SNAP/res/sample-configuration.toml" "$HOME/.edgex-cli/configuration.toml"
    echo "Created $HOME/.edgex-cli/configuration.toml from sample configuration file."
fi

exec "$@"
