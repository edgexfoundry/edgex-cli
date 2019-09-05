package version

import (
	"fmt"

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
			fmt.Println("edgex-cli version: ", version)
		},
	}

	return cmd
}
