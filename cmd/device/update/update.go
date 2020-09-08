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

const deviceTemp = `{
  "Id": "{{.Id}}",
  "Name": "{{.Name}}",
  "Description":"{{.Description}}",
  "Adminstate":"{{.AdminState}}",
  "OperatingState": "{{.OperatingState}}",
  "Protocols":{ 
     {{- $length := len .Protocols }}{{$idx:=0}}
     {{- range $k,$ProtocolProperties := .Protocols}}
     "{{$k}}":{ 
	    {{- $innerLength := len $ProtocolProperties }}{{$innerIdx:=0}} 
        {{- range $ki,$vi := $ProtocolProperties}}
          "{{- $ki}}":"{{$vi}}"{{if not (lastElem $innerIdx $innerLength)}},{{end}}{{$innerIdx =  inc $innerIdx}}
        {{end}}} {{if not (lastElem $idx $length)}},{{end}}{{$idx =  inc $idx}} 
     {{end}}}, 
  "Service": {
     "Name": "{{.Service.Name}}"
  },
  "Profile": {
    "Name": "{{.Profile.Name}}"
  }
}`

var name string
var file string

// NewCommand returns the update device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update device",
		Long: `Update device with given name using interactive mode or update device(s) described in a JSON file.
Interactive mode opens a default editor to customize Device information`,
		RunE: deviceHandler,
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Json file containing device service configuration to update")
	return cmd
}

func deviceHandler(cmd *cobra.Command, args []string) error {
	if name != "" && file != "" {
		return errors.New("you could work with interactive mode by providing the name of the device you want to update or by providing a file, but not with both")
	}

	if name == "" && file == "" {
		return errors.New("Please, provide file or device name ")
	}

	if file != "" {
		return updateDevicesFromFile()
	}
	//Update the device provided by name using interactive mode to alter it
	d, err := parseDevice(name)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
	err = metadata.NewDeviceClient(client).Update(context.Background(), d)
	if err != nil {
		return err
	}

	return nil
}

func updateDevicesFromFile() error {
	devices, err := LoadDevicesFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
	for _, d := range devices {
		err = metadata.NewDeviceClient(client).Update(context.Background(), d)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

//loads a device service to be updated and open a default editor for customization
func parseDevice(name string) (models.Device, error) {
	//load Device from database
	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
	device, err := metadata.NewDeviceClient(client).DeviceForName(context.Background(), name)
	if err != nil {
		return models.Device{}, err
	}

	//populate the template with the loaded device and open default editor, so the client could customize the data
	updatedDeviceBytes, err := editor.OpenInteractiveEditor(device, deviceTemp, template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"lastElem": editor.IsLastElementOfSlice,
	})
	if err != nil {
		return models.Device{}, err
	}

	var d models.Device
	err = json.Unmarshal(updatedDeviceBytes, &d)
	if err != nil {
		return models.Device{}, errors.New("Unable to execute the command. The provided information is not valid: " + err.Error())
	}
	return d, nil
}


//LoadDevicesFromFile could read a file that contains single Device or list of Device
func LoadDevicesFromFile(filePath string) ([]models.Device, error){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var devices []models.Device

	//check if the file contains just one Device
	var d models.Device
	err = json.Unmarshal(file, &d)
	if err != nil {
		//check if the file contains list of Device
		err = json.Unmarshal(file, &devices)
		if err != nil {
			return nil, err
		}
	} else {
		devices = append(devices, d)
	}
	return devices, nil
}

