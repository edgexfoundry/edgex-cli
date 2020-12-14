package adminstate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/go-mod-core-contracts/requests/states/admin"

	"github.com/spf13/cobra"
)

var deviceName string
var deviceId string
var state string

// NewCommand returns device adminstate
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adminstate",
		Short: "Update deviceName admin state",
		Long:  `Update deviceName admin state`,
		RunE:  updateAdminStatusHandler,
	}
	cmd.Flags().StringVarP(&state, "state", "s", "", fmt.Sprintf("Admin State \nValues: [%s]", strings.Join([]string{models.Locked, models.Unlocked}, ",")))
	cmd.Flags().StringVarP(&deviceName, "name", "n", "", "Name of the deviceName to be updated")
	cmd.Flags().StringVarP(&deviceId, "id", "i", "", "Id of the deviceName to be updated")
	return cmd
}

func updateAdminStatusHandler(cmd *cobra.Command, args []string) error {
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
		err = mdc.UpdateAdminStateByName(cmd.Context(), deviceName, admin.UpdateRequest{AdminState: models.AdminState(state)})
	} else {
		err = mdc.UpdateAdminState(cmd.Context(), deviceId, admin.UpdateRequest{AdminState: models.AdminState(state)})
	}
	return err
}

func validate() error {
	if state == "" {
		return errors.New("admin state should be specified")
	} else if state != "" {
		valid, err := models.AdminState(state).Validate()
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
