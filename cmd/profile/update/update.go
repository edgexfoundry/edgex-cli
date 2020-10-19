package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/edgex-cli/cmd/profile/add"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
	"html/template"
)

var file string
var name string
var description string
var manufacturer string
var model string
var labels string

// NewCommand returns the update device profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update device service",
		Long: `Update device profile(s) described in the given JSON/YAML file or use the interactive mode enabled by providing 
 name of existing device profile`,
		RunE: handler,
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVar(&manufacturer, "manufacturer", "", "Manufacturer")
	cmd.Flags().StringVar(&model, "model", "", "Model")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma separated strings")

	cmd.Flags().StringVarP(&file, "file", "f", "", "Json/YAML file containing device service configuration to update")
	return cmd
}

func handler(cmd *cobra.Command, args []string) error {
	if name != "" && file != "" {
		return errors.New("Profile could be updated by providing a file, or by specifying device service name to be updated using interactive mode. ")
	}

	if name == "" && file == "" {
		return errors.New("Please, provide file or profile name ")
	}

	if file != "" {
		return updateProfileFromFile()
	}

	updatedProfile, err := parseProfile(name)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	err = metadata.NewDeviceProfileClient(client).Update(context.Background(), updatedProfile)
	if err != nil {
		return err
	}

	return nil
}

//parseProfile loads a DeviceProfile to be updated and open a default editor for customization
func parseProfile(name string) (models.DeviceProfile, error) {
	var err error
	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	ds, err := metadata.NewDeviceProfileClient(client).DeviceProfileForName(context.Background(), name)
	if err != nil {
		return models.DeviceProfile{}, err
	}
	updatedProfileBytes, err := editor.OpenInteractiveEditor(ds, add.ProfileTemplate, template.FuncMap{
		"lastElem":   editor.IsLastElementOfSlice,
		"EscapeHTML": editor.EscapeHTML,
	})

	if err != nil {
		return models.DeviceProfile{}, err
	}
	var updatedProfile models.DeviceProfile
	err = json.Unmarshal(updatedProfileBytes, &updatedProfile)
	if err != nil {
		return models.DeviceProfile{}, errors.New("Unable to execute the command. The provided information is invalid: " + err.Error())
	}
	return updatedProfile, err
}

func updateProfileFromFile() error {
	profiles, err := add.LoadFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	for _, ds := range profiles {
		err = metadata.NewDeviceProfileClient(client).Update(context.Background(), ds)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}
