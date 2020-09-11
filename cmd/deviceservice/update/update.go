package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

const DeviceServiceTempl = `{
   "Id": "{{.Id}}",
   "Name" : "{{.Name}}",
   "Description" : "{{.Description}}",
   "AdminState" : "{{.AdminState}}",
   "OperatingState" : "{{.OperatingState}}",
   "Labels" : [
    {{- $labelsLenght := len .Labels}}
    {{- range $idx, $l := .Labels}}
      "{{$l}}" {{if not (lastElem $idx $labelsLenght)}},{{end}} 
    {{- end}}
    ],
   "Addressable" :
   {
     "Name" : "{{.Addressable.Name}}"
   }
 }`

var name string
var file string

// NewCommand returns the update device service command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update device service",
		Long:  `Update device service(s) described in the given JSON file or use the interactive mode enabled by providing 
 name of existing device service.`,
		RunE:  deviceServiceHandler,
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Device Service name. Service with given name is loaded into default editor, ready to be customized")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Json file containing device service configuration to update")
	return cmd
}

func deviceServiceHandler(cmd *cobra.Command, args []string) error {
	if name != "" && file != "" {
		return errors.New("DeviceService could be updated by providing a file, or by specifying device service name to be updated using interactive mode. ")
	}

	if name == "" && file == "" {
		return errors.New("Please, provide file or device service name ")
	}

	if file != "" {
		return updateDeviceServiceFromFile()
	}

	updatedDeviceService, err := parseDeviceService(name)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	err = metadata.NewDeviceServiceClient(client).Update(context.Background(), updatedDeviceService)
	if err != nil {
		return err
	}

	return nil
}
//parseDeviceService loads a device service to be updated and open a default editor for customization
func parseDeviceService(name string) (models.DeviceService, error) {
	var err error
	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	ds, err := metadata.NewDeviceServiceClient(client).DeviceServiceForName(context.Background(), name)
	if err != nil {
		return models.DeviceService{}, err
	}

	updatedDeviceServiceBytes, err := editor.OpenInteractiveEditor(ds, DeviceServiceTempl, template.FuncMap{
		"lastElem": editor.IsLastElementOfSlice,
	})

	if err != nil {
		return models.DeviceService{}, err
	}
	var updatedDeviceService models.DeviceService
	err = json.Unmarshal(updatedDeviceServiceBytes, &updatedDeviceService)
	if err != nil {
		return models.DeviceService{}, errors.New("Unable to execute the command. The provided information is not valid: " + err.Error())
	}
	return updatedDeviceService, err
}

func updateDeviceServiceFromFile() error {
	deviceServices, err := LoadDSFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	for _, ds := range deviceServices {
		err = metadata.NewDeviceServiceClient(client).Update(context.Background(), ds)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

//loadJsonFile could read a file that contains single Device service or list of Device Services
func LoadDSFromFile(filePath string) ([]models.DeviceService, error){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var deviceServices []models.DeviceService

	//check if the file contains just one DeviceService
	var ds models.DeviceService
	err = json.Unmarshal(file, &ds)
	if err != nil {
		//check if the file contains list of DeviceServices
		err = json.Unmarshal(file, &deviceServices)
		if err != nil {
			return nil, err
		}
	} else {
		deviceServices = append(deviceServices, ds)
	}
	return deviceServices, nil
}
