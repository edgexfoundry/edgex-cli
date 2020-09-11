package add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"github.com/edgexfoundry-holding/edgex-cli/cmd/deviceservice/update"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var interactiveMode bool
var name string
var description string
var adminState string
var operState string
var addrName string
var file string

const DeviceServicesTempl = `[{{range $ds := .}}` + update.DeviceServiceTempl +
	`{{end}}]`

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add device service",
		Long:  `Create device service(s) described in the given JSON file or use the interactive mode with additional flags.`,
		RunE:  newDeviceServiceHandler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Event information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVar(&adminState, "adminState", "", "Admin Status")
	cmd.Flags().StringVar(&operState, "operatingStatus", "", "Operating Status")
	cmd.Flags().StringVar(&addrName, "addrName", "", "Addressable name the device service will be linked with.")

	cmd.Flags().StringVarP(&file, "file", "f", "", "Json file containing device service(s) configuration")
	return cmd
}

func newDeviceServiceHandler(cmd *cobra.Command, args []string) error {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	if file != "" {
		return createDeviceServicesFromFile()
	}

	deviceServices, err := parseDeviceService(interactiveMode)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	for _, ds := range deviceServices {
		_, err = metadata.NewDeviceServiceClient(client).Add(context.Background(), &ds)
		if err != nil {
			return err
		}
	}
	return nil
}

func populateDeviceService(deviceServices *[]models.DeviceService) {
	ds := models.DeviceService{}
	ds.Name = name
	ds.Description = description
	ds.AdminState = models.AdminState(adminState)
	ds.OperatingState = models.OperatingState(operState)
	ds.Addressable = models.Addressable{Name: addrName}
	*deviceServices = append(*deviceServices, ds)
}

func createDeviceServicesFromFile() error {
	deviceServices, err := update.LoadDSFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	for _, ds := range deviceServices {
		_, err = metadata.NewDeviceServiceClient(client).Add(context.Background(), &ds)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

func parseDeviceService(interactiveMode bool) ([]models.DeviceService, error) {
	//parse Device Service based on interactive mode and the other provided flags
	var err error
	var deviceServices []models.DeviceService
	populateDeviceService(&deviceServices)

	var updatedDeviceServiceBytes []byte
	if interactiveMode {
		updatedDeviceServiceBytes, err = editor.OpenInteractiveEditor(deviceServices, DeviceServicesTempl, template.FuncMap{
			"lastElem": editor.IsLastElementOfSlice,
		})
	} else {
		updatedDeviceServiceBytes, err = json.Marshal(deviceServices)
	}
	if err != nil {
		return nil, err
	}

	var updatedDeviceServices []models.DeviceService
	err = json.Unmarshal(updatedDeviceServiceBytes, &updatedDeviceServices)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is not valid: " + err.Error())
	}
	return updatedDeviceServices, err
}
