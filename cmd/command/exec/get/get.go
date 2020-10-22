package get

import (
	"context"
	"fmt"
	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/command"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

var deviceName string
var cmdName string
var data string

// NewCommand returns `issue Get command` command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Issue GET command",
		Long:  `Issue GET command referenced by the command name to the device also referenced by name`,
		RunE:  handler,
	}
	cmd.Flags().StringVarP(&deviceName, "device", "d", "", "Specify the name of device")
	cmd.Flags().StringVarP(&cmdName, "cmd", "n", "", "Specify the name of the command to be executed")
	cmd.MarkFlagRequired("device")
	cmd.MarkFlagRequired("cmd")
	return cmd
}

func handler(cmd *cobra.Command, args []string) error {
	client := local.New(config.Conf.Clients["Command"].Url() + clients.ApiDeviceRoute)
	res, err := command.NewCommandClient(client).GetDeviceCommandByNames(context.Background(), deviceName, cmdName)
	if err == nil {
		fmt.Printf("%s\n", res)
	}
	return err
}
