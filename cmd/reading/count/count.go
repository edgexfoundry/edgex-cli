package count

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/formatters"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/coredata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

// NewCommand returns the count reading command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Returns the count of core-data readings",
		Long:  `Return a count of the number of readings in core data.`,
		Args:  cobra.MaximumNArgs(0),
		RunE:  countHandler,
	}
	return cmd
}

func countHandler(cmd *cobra.Command, args []string) (err error) {
	client := coredata.NewReadingClient(
		local.New(config.Conf.Clients["CoreData"].Url() + clients.ApiReadingRoute),
	)
	count, err := client.ReadingCount(context.Background())
	if err != nil {
		return
	}

	formatter := formatters.NewFormatter(fmt.Sprintf("Total readings count: %v", count), nil)
	err = formatter.Write(count)
	return
}
