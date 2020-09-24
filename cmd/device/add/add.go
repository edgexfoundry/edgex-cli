// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"io/ioutil"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/spf13/cobra"
)

const DeviceTemplate = `[{{range $d := .}}` + update.DeviceTemp + `{{end}}]`


var interactiveMode bool
var name string
var description string
var adminState string
var operState string
var profileName string
var serviceName string
var file string

// NewCommand returns the update device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add devices",
		Long:  `Create devices described in a given JSON file or use the interactive mode with additional flags.`,
		RunE:  deviceHandler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Event information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVar(&adminState, "adminState", "", "Admin Status")
	cmd.Flags().StringVar(&operState, "operatingStatus", "", "Operating Status")
	cmd.Flags().StringVar(&profileName, "profileName", "", "Device Profile name")
	cmd.Flags().StringVar(&serviceName, "serviceName", "", "Device Service name")

	cmd.Flags().StringVarP(&file, "file", "f", "", "File containing device configuration in json format")
	return cmd
}

func deviceHandler(cmd *cobra.Command, args []string) error {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	if file != "" {
		return createDeviceFromFile()
	}

	interactiveMode, err := cmd.Flags().GetBool(editor.InteractiveModeLabel)
	if err != nil {
		return err
	}

	devices, err := parseDevice(interactiveMode)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
	for _, d := range devices {
		_, err = metadata.NewDeviceClient(client).Add(context.Background(), &d)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDeviceFromFile() error {
	devices, err := LoadDeviceFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
	for _, d := range devices {
		_, err = metadata.NewDeviceClient(client).Add(context.Background(), &d)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

func parseDevice(interactiveMode bool) ([]models.Device, error) {
	//parse Device based on interactive mode and the other provided flags
	var err error
	var devices []models.Device
	if name != "" {
		client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute)
		from, _ := metadata.NewDeviceClient(client).DeviceForName(context.Background(), name)
		devices = append(devices, from)
	} else {
		populateDevice(&devices)
	}

	var updatedDeviceBytes []byte
	if interactiveMode {
		updatedDeviceBytes, err = editor.OpenInteractiveEditor(devices, DeviceTemplate, nil)
	}
	if err != nil {
		return nil, err
	}

	var updatedDevices []models.Device
	err = json.Unmarshal(updatedDeviceBytes, &updatedDevices)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is not valid:" + err.Error())
	}
	return updatedDevices, err
}

func populateDevice(devices *[]models.Device) {
	d := models.Device{}
	d.Name = name
	d.Description = description
	d.AdminState = models.AdminState(adminState)
	d.OperatingState = models.OperatingState(operState)
	d.Profile = models.DeviceProfile{Name: profileName}
	d.Service = models.DeviceService{Name: serviceName}
	*devices = append(*devices, d)
}

//LoadDeviceFromFile could read a file that contains single Device or list of Device Services
func LoadDeviceFromFile(filePath string) ([]models.Device, error) {
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
