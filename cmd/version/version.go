package version

import (
	"fmt"

	client "github.com/edgexfoundry/edgex-cli/pkg"
	"github.com/spf13/cobra"
)

var version string

// NewCommand returns the version command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Version command",
		Long:  `Outputs the current version of edgex-cli.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Printing the CLI's version
			fmt.Println("edgex-cli version: ", version)
			// printing the version of each service
			// ports :=
			data := client.GetVersion("48080")
			fmt.Println("version: ", string(data))
		},
	}

	return cmd
}
