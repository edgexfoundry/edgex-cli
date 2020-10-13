package operstatus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/go-mod-core-contracts/requests/states/operating"

	"github.com/spf13/cobra"
)

var deviceName string
var deviceId string
var state string

// NewCommand returns device operstatus
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operstate",
		Short: "Update deviceName operating state",
		Long:  `Update deviceName operating state`,
		RunE:  updateOperStatusHandler,
	}
	cmd.Flags().StringVarP(&state, "state", "s", "", fmt.Sprintf("Operating State \nValues: [%s]", strings.Join([]string{models.Disabled, models.Enabled}, ",")))
	cmd.Flags().StringVarP(&deviceName, "name", "n", "", "Name of the deviceName to be updated")
	cmd.Flags().StringVarP(&deviceId, "id", "i", "", "Id of the deviceName to be updated")
	return cmd
}

func updateOperStatusHandler(cmd *cobra.Command, args []string) error {
	state = strings.ToUpper(state)
	if err := validate(); err != nil {
		return err
	}
	url := config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute
	mdc := metadata.NewDeviceClient(
		local.New(url),
	)
	var err error
	if deviceName != "" {
		err = mdc.UpdateOpStateByName(context.Background(), deviceName, operating.UpdateRequest{OperatingState: models.OperatingState(state)})
	} else {
		err = mdc.UpdateOpState(context.Background(), deviceId, operating.UpdateRequest{OperatingState: models.OperatingState(state)})
	}
	return err
}

func validate() error {
	if state == "" {
		return errors.New("operating state should be specified")
	} else if state != "" {
		valid, err := models.OperatingState(state).Validate()
		if !valid {
			return err
		}
	} else if deviceName == "" && deviceId == "" {
		return errors.New("deviceName or deviceId should be specified")
	} else if deviceName != "" && deviceId != "" {
		return errors.New("only one deviceName or deviceId should be specified")
	}
	return nil
}
