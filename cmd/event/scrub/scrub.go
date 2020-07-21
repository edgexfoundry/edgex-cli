package scrub

import (
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/confirmation"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

var all bool

// NewCommand return scrub events command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "scrub",
		Short: "Remove all (pushed) events and their associated readings [USE WITH CAUTION]",
		Long: `[USE WITH CAUTION] The effect of this command is irreversible. 
It removes all pushed events and their associated.
When used with "--all" flag, it removes all readings and events from the database`,
		RunE: scrubHandler,
	}
	cmd.Flags().BoolP("all", "a", false, "Removes all readings and events from the database [USE WITH CAUTION]")
	return cmd
}

func scrubHandler(cmd *cobra.Command, args []string) (err error) {
	// asking user to confirm the scrub command
	if !confirmation.New().Confirm() {
		return
	}
	all, err1 := cmd.Flags().GetBool("all")
	if err1 != nil {
		return err1
	}

	url := config.Conf.Clients["CoreData"].Url() + clients.ApiEventRoute + "/scrub"
	if all {
		url = config.Conf.Clients["CoreData"].Url() + clients.ApiEventRoute + "/scruball"
	}

	err = request.Delete(url)
	if err != nil {
		return err
	}

	return
}
