package adminstate

import (
	"github.com/edgexfoundry/edgex-cli/cmd/device/adminstate/update"

	"github.com/spf13/cobra"
)

// NewCommand returns device adminstate
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adminstate",
		Short: "Device admin state ",
	}
	cmd.AddCommand(update.NewCommand())
	return cmd
}
