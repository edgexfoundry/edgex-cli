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
	jsonpkg "encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
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
	initNameDeviceProfileCommand(cmd)
}

func initRmDeviceProfileCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a device profile",
		Long:  "Removes a device profile from the core-metadata database",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			response, err := getCoreMetaDataService().RemoveDeviceProfile(deviceProfileName)
			if err == nil && response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&deviceProfileName, "name", "n", "", "Device Profile name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

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

func initNameDeviceProfileCommand(cmd *cobra.Command) {
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

func handleGetDeviceProfileByName(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().GetDeviceProfileByNameJSON(deviceProfileName)
		if err == nil {
			fmt.Print(json)
		}
		return err
	} else {
		p, err := getCoreMetaDataService().GetDeviceProfileByName(deviceProfileName)
		if p != nil {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			printProfileTableHeader(w)
			printProfile(w, p)
			w.Flush()
		}
		return err
	}
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

	resources, commands, labels, err := getDeviceProfileAttributes()
	if err != nil {
		return err
	}

	response, err := getCoreMetaDataService().AddDeviceProfile(deviceProfileName, deviceProfileDescription, deviceProfileManufacturer,
		deviceProfileModel, labels, resources, commands)

	if err == nil && response != nil {
		fmt.Println(response)
	}
	return err
}

func handleListDeviceProfile(cmd *cobra.Command, args []string) error {

	if json {
		json, _, err := getCoreMetaDataService().ListAllDeviceProfilesJSON(offset, limit, labels)
		if err != nil {
			return err
		}
		fmt.Print(json)

	} else {

		deviceProfiles, err := getCoreMetaDataService().ListAllDeviceProfiles(offset, limit, getLabels())
		if err != nil {
			return err
		}
		if len(deviceProfiles) == 0 {
			fmt.Println("No device profiles available")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printProfileTableHeader(w)
		for _, p := range deviceProfiles {
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
			time.Unix(0, p.Created).Format(time.RFC822),
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
