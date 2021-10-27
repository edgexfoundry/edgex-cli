/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package cmd

import (
	"context"
	jsonpkg "encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/spf13/cobra"
)

var deviceProfileName, deviceProfileResources, deviceProfileDescription string
var deviceProfileManufacturer, deviceProfileModel, deviceProfileCommands string

func init() {
	var cmd = &cobra.Command{
		Use:   "deviceprofile",
		Short: "Add, remove, get and list device profiles [Core Metadata]",
		Long:  "Add, remove, get and list device profiles [Core Metadata]",

		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)

	initRmDeviceProfileCommand(cmd)
	initListDeviceProfileCommand(cmd)
	initAddDeviceProfileCommand(cmd)
	initGetDeviceProfileByNameCommand(cmd)
}

// initRmDeviceProfileCommand implements the DELETE ​/device​profile/name​/{name} endpoint
// "Delete a device profile by its unique name. This operation will fail if there are devices actively using the profile."
func initRmDeviceProfileCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:          "rm",
		Short:        "Remove a device profile",
		Long:         "Removes a device profile from the core-metadata database",
		RunE:         handleRmDeviceProfile,
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&deviceProfileName, "name", "n", "", "Device Profile name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

// initAddDeviceProfileCommand implements the POST ​/device​profile endpoint
// "Allows creation of a new device profile"
func initAddDeviceProfileCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "add",
		Short: "Add a new device profile",
		Long: `Add a new device profile.
 
 Example:
 ./bin/edgex-cli deviceprofile add 
	 -n testprofile 
	 -r "[{\"name\": \"SwitchButton\",\"description\": \"Switch On/Off.\",\"properties\": {\"valueType\": \"String\",\"readWrite\": \"RW\",\"defaultValue\": \"On\",\"units\": \"On/Off\" } }]" 
	 -c "[{\"name\": \"Switch\",\"readWrite\": \"RW\",\"resourceOperations\": [{\"deviceResource\": \"SwitchButton\",\"DefaultValue\": \"false\" }]} ]"	 
 `,
		RunE:         handleAddDeviceProfile,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&deviceProfileName, "name", "n", "", "Device profile name")
	add.Flags().StringVarP(&deviceProfileDescription, "description", "d", "", "Device profile description")
	add.Flags().StringVarP(&deviceProfileManufacturer, "manufacturer", "m", "", "Manufacturer of the device")
	add.Flags().StringVarP(&deviceProfileModel, "model", "", "", "Model of the device")
	add.Flags().StringVarP(&deviceProfileResources, "resources", "r", "", "JSON structure representing a device resource that can be read or written")
	add.Flags().StringVarP(&deviceProfileCommands, "commands", "c", "", "JSON structure defining read/write capabilities native to the device")
	add.MarkFlagRequired("name")
	addLabelsFlag(add)
	cmd.AddCommand(add)
}

// initListDeviceProfileCommand implements the GET ​/device​profile/all endpoint
// "Given the entire range of device profiles sorted by last modified descending, returns a portion
// of that range according to the offset and limit parameters. Device profiles may also be filtered by label."
func initListDeviceProfileCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List device profiles",
		Long:         `List all device profiles, optionally specifying a limit, offset and/or label(s)`,
		RunE:         handleListDeviceProfile,
		SilenceUsage: true,
	}

	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	addLabelsFlag(listCmd)
}

// initGetDeviceProfileByNameCommand implements the GET ​/device​profile/name endpoint
// "Returns a device profile by its name"
func initGetDeviceProfileByNameCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Returns a device profile by name",
		Long:         `Returns a device profile by name`,
		RunE:         handleGetDeviceProfileByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&deviceProfileName, "name", "n", "", "Device profile name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleRmDeviceProfile(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceProfileClient()
	response, err := client.DeleteByName(context.Background(), deviceProfileName)
	if err == nil {
		fmt.Println(response)
	}
	return err
}

func handleGetDeviceProfileByName(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceProfileClient()

	response, err := client.DeviceProfileByName(context.Background(), deviceProfileName)
	if err != nil {
		return err
	}

	if json {
		result, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Println(string(result))
	} else {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printProfileTableHeader(w)
		printProfile(w, &response.Profile)
		w.Flush()
	}
	return nil

}

func getDeviceProfileAttributes() (resources []dtos.DeviceResource, commands []dtos.DeviceCommand, labels []string, err error) {

	labels = getLabels()
	if deviceProfileResources != "" {
		err = jsonpkg.Unmarshal([]byte(deviceProfileResources), &resources)
		if err != nil {
			err = errors.New(`please specify device resources using a JSON value. Example: -r "[{\"name\": \"SwitchButton\",\"description\": \"Switch On/Off.\",\"properties\": {\"valueType\": \"String\",\"readWrite\": \"RW\",\"defaultValue\": \"On\",\"units\": \"On/Off\" } }]"`)
		}
	}
	if deviceProfileCommands != "" {
		err = jsonpkg.Unmarshal([]byte(deviceProfileCommands), &commands)
		if err != nil {
			err = errors.New(`please specify device resources using a JSON value. Example: -c "[{\"name\": \"Switch\",\"readWrite\": \"RW\",\"resourceOperations\": [{\"deviceResource\": \"SwitchButton\",\"DefaultValue\": \"false\" }]} ]"`)
		}
	}
	return
}

func handleAddDeviceProfile(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceProfileClient()

	resources, commands, labels, err := getDeviceProfileAttributes()
	if err != nil {
		return err
	}

	var req = requests.NewDeviceProfileRequest(dtos.DeviceProfile{
		Name:            deviceProfileName,
		Description:     deviceProfileDescription,
		Manufacturer:    deviceProfileManufacturer,
		Model:           deviceProfileModel,
		Labels:          labels,
		DeviceResources: resources,
		DeviceCommands:  commands,
	})

	response, err := client.Add(context.Background(), []requests.DeviceProfileRequest{req})

	if response != nil {
		fmt.Println(response[0])
	}
	return err

}

func handleListDeviceProfile(cmd *cobra.Command, args []string) error {

	client := getCoreMetaDataService().GetDeviceProfileClient()
	response, err := client.AllDeviceProfiles(context.Background(), getLabels(), offset, limit)
	if err != nil {
		return err
	}
	if json {
		result, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Println(string(result))
	} else {

		if len(response.Profiles) == 0 {
			fmt.Println("No profiles available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printProfileTableHeader(w)
		for _, p := range response.Profiles {
			printProfile(w, &p)
		}
		w.Flush()
	}
	return nil
}

func printProfileTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tCreated\tDescription\t# DeviceCommands\t# DeviceResources\tManufacturer\tModel\tName")

	} else {
		fmt.Fprintln(w, "Name\tDescription\tManufacturer\tModel\tName")
	}

}

func printProfile(w *tabwriter.Writer, p *dtos.DeviceProfile) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			p.Id,
			p.Name,
			getRFC822Time(p.Created),
			p.Description,
			len(p.DeviceCommands),
			len(p.DeviceResources),
			p.Manufacturer,
			p.Model,
			p.Name)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
			p.Name,
			p.Description,
			p.Manufacturer,
			p.Model,
			p.Name)
	}

}
