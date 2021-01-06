package put

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/command"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

var device string
var cmdName string
var body string
var file string

// NewCommand returns `issue Put command` command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "put",
		Short: "Issue PUT command",
		Long:  `Issue PUT command referenced by the command name to the device also referenced by name`,
		RunE:  handler,
	}
	cmd.Flags().StringVarP(&device, "device", "d", "", "Specify the name of the device")
	cmd.Flags().StringVarP(&cmdName, "cmd", "n", "", "Specify the name of the command to be executed")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Specify PUT requests body/data inline")
	cmd.Flags().StringVarP(&file, "file", "f", "", "File containing PUT requests body/data")
	cmd.MarkFlagRequired("device")
	cmd.MarkFlagRequired("cmdName")
	return cmd
}

func handler(cmd *cobra.Command, args []string) (err error) {
	if (file != "" && body != "") || (file == "" && body == "") {
		return errors.New("please specify request data using one of the provided ways: --body or --file  ")
	}
	if file != "" {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		body = string(content)
	}

	client := local.New(config.Conf.Clients["Command"].Url() + clients.ApiDeviceRoute)
	res, err := command.NewCommandClient(client).PutDeviceCommandByNames(cmd.Context(), device, cmdName, body)
	if res != "" {
		fmt.Printf("%s\n", res)
	}
	return err
}
