package add

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/editor"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"html/template"
	"io/ioutil"

	"github.com/spf13/cobra"
)

const dsTemplate = `{
  "Name" : "{{.Name}}",
  "Description" : "{{.Description}}",
  "AdminState" : "{{.AdminState}}",
  "OperatingState" : "{{.OperatingState}}",
  "Labels" : [],
  "Addressable" : 
    {
      "Name" : "{{.Addressable.Name}}"
    }
}`

var interactiveMode bool
var name string
var description string
var adminState string
var operState string
var addrName string
var file string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add device service",
		Long:  `Create device service described in the given JSON file or use the interactive mode with additional flags.`,
		RunE: newDeviceServiceHandler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Event information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVar(&adminState, "adminState","", "Admin Status")
	cmd.Flags().StringVar(&operState, "operatingStatus","", "Operating Status")
	cmd.Flags().StringVar(&addrName, "addrName","", "Addressable name the device service will be linked with.")

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

    //create Device Service based on interactive mode and the other provided flags
	var err error
	ds := models.DeviceService{}
	populateDeviceService(&ds)

	var updatedDeviceServiceBytes []byte
	interactiveMode, err := cmd.Flags().GetBool(editor.InteractiveModeLabel)
	if interactiveMode {
		updatedDeviceServiceBytes, err = openInteractiveEditor(ds)
	} else {
		updatedDeviceServiceBytes, err = json.Marshal(ds)
	}
	if err != nil {
		return err
	}

	updatedDeviceService := models.DeviceService{}
	err = json.Unmarshal(updatedDeviceServiceBytes, &updatedDeviceService)
	if err != nil {
		return errors.New("Unable to create a Device Service. The provided information is not valid" + err.Error())
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute)
	_, err = metadata.NewDeviceServiceClient(client).Add(context.Background(), &updatedDeviceService)
	if err != nil {
		return err
	}
	return nil
}

func populateDeviceService(ds *models.DeviceService) {
	ds.Name = name
	ds.Description = description
	ds.AdminState = models.AdminState(adminState)
	ds.OperatingState = models.OperatingState(operState)
	ds.Addressable = models.Addressable{Name: addrName}
}

// openInteractiveEditor opens the users default editor and with a JSON representation of the Device Service.
func openInteractiveEditor(ds models.DeviceService) ([]byte, error) {
	dsJsonTemplate, err := template.New("DS").Parse(dsTemplate)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	err = dsJsonTemplate.Execute(buff, ds)
	if err != nil {
		return nil, err
	}

	return editor.CaptureInputFromEditor(buff.Bytes())
}

func createDeviceServicesFromFile() error {
	deviceServices, err := loadJsonFile()
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
//loadJsonFile could read a file that contains single Device service or list of Device Services
func loadJsonFile() ([]models.DeviceService, error){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(file)
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