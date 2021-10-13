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

func init() {
	var cmd = &cobra.Command{
		Use:          "device",
		Short:        "Add, remove, get, list and modify devices [Core Metadata]",
		Long:         "Add, remove, get, list and modify devices [Core Metadata]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initAddDeviceCommand(cmd)
	initListDeviceCommand(cmd)
	initNameDeviceCommand(cmd)
	initRmDeviceCommand(cmd)
	initUpdateDeviceCommand(cmd)

}

var deviceName, deviceId, deviceDescription, deviceAdminState, deviceOperState, deviceProfile, deviceService string
var deviceLocation, deviceProtocols string

// initRmDeviceCommand implements the DELETE ​/device​/name​/{name} endpoint
// "Delete a device by name"
func initRmDeviceCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a device",
		Long:  "Removes a device from the core-metadata database",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := getCoreMetaDataService().RemoveDevice(deviceName)
			if err == nil && response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&deviceName, "name", "n", "", "Device name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

// initListDeviceCommand implements the GET ​/device​/all endpoint:
// "Given the entire range of devices sorted by last modified descending,
// returns a portion of that range according to the offset and limit parameters.
// Devices may also be filtered by label."
func initListDeviceCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List devices",
		Long:         `List all devices, optionally specifying a limit, offset and/or label(s)`,
		RunE:         handleListDevices,
		SilenceUsage: true,
	}

	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	addLabelsFlag(listCmd)
	cmd.AddCommand(listCmd)
}

// initUpdateDeviceCommand implements the PATCH ​/device endpoint
// "Allows updates to an existing device"
func initUpdateDeviceCommand(cmd *cobra.Command) {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing device",
		Long: `Update an existing device 
'id' and 'name' must be specified in order to identify the service. 
Any other provided non-blank property will be updated.

Example: 
 edgex-cli device update -n AWS IOT Button1 -i "edaa7c0f-05c6-4368-89f1-3be5e197cf6a" -l "new-label"
		`,
		RunE:         handleUpdateDevice,
		SilenceUsage: true,
	}
	updateCmd.Flags().StringVarP(&deviceName, "name", "n", "", "Device name")
	updateCmd.Flags().StringVarP(&deviceId, "id", "i", "", "Device name")
	updateCmd.Flags().StringVarP(&deviceDescription, "description", "d", "", "Device description")
	updateCmd.Flags().StringVarP(&deviceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	updateCmd.Flags().StringVarP(&deviceOperState, "operating-state", "o", "UP", "Operating state [UP | DOWN | UNKNOWN]")
	updateCmd.Flags().StringVarP(&deviceProfile, "profile", "p", "", "Associated device profile")
	updateCmd.Flags().StringVarP(&deviceService, "service", "s", "", "Associated device service")
	updateCmd.Flags().StringVarP(&deviceLocation, "location", "l", "", "Device location")
	updateCmd.Flags().StringVarP(&deviceProtocols, "protocols", "", "", "A map of supported protocols")
	addLabelsFlag(updateCmd)
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("id")
	cmd.AddCommand(updateCmd)
}

// initAddDeviceCommand implements the POST ​/device endpoint
// "Allows provisioning of a new device"
func initAddDeviceCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "add",
		Short: "Provision a new device",
		Long: `Provision a new device

Example: 
	edgex-cli device add -n TestDevice -p TestDeviceProfile -s TestDeviceService --protocols "{\"modbus-tcp\":{\"Address\": \"localhost\",\"Port\": \"1234\" }}"

		
`,
		RunE:         handleAddDevice,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&deviceName, "name", "n", "", "Device name")
	add.Flags().StringVarP(&deviceDescription, "description", "d", "", "Device description")
	add.Flags().StringVarP(&deviceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	add.Flags().StringVarP(&deviceOperState, "operating-state", "o", "UP", "Operating state [UP | DOWN | UNKNOWN]")
	add.Flags().StringVarP(&deviceProfile, "profile", "p", "", "Associated device profile")
	add.Flags().StringVarP(&deviceService, "service", "s", "", "Associated device service")
	add.Flags().StringVarP(&deviceLocation, "location", "l", "", "Device location")
	add.Flags().StringVarP(&deviceProtocols, "protocols", "", "", "A map of supported protocols")
	addLabelsFlag(add)
	add.MarkFlagRequired("name")
	add.MarkFlagRequired("service")
	add.MarkFlagRequired("profile")
	add.MarkFlagRequired("protocols")
	cmd.AddCommand(add)
}

func initNameDeviceCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Returns a device by name",
		Long:         `Returns a device by name`,
		RunE:         handleGetDeviceByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&deviceName, "name", "n", "", "Device name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleUpdateDevice(cmd *cobra.Command, args []string) error {
	var name, id, description, service, profile, admin, oper, location *string
	if deviceName != "" {
		name = &deviceName
	}
	if deviceId != "" {
		id = &deviceId
	}
	if deviceDescription != "" {
		description = &deviceDescription
	}

	if deviceService != "" {
		service = &deviceService
	}

	if deviceProfile != "" {
		profile = &deviceProfile
	}

	if deviceAdminState != "" {
		admin = &deviceAdminState
	}

	if deviceOperState != "" {
		oper = &deviceOperState
	}

	if deviceLocation != "" {
		location = &deviceLocation
	}

	protocols, labels, err := getDeviceAttributes()
	if err != nil {
		return err
	}

	response, err := getCoreMetaDataService().UpdateDevice(name, id, description, service, profile, admin, oper, location, labels, protocols)
	if response != nil {
		fmt.Println(response)
	}
	return err

}

func handleGetDeviceByName(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().GetDeviceByNameJSON(deviceName)
		if err == nil {
			fmt.Print(json)
		}
		return err
	} else {
		d, err := getCoreMetaDataService().GetDeviceByName(deviceName)
		if d != nil {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			printDeviceTableHeader(w)
			printDevice(w, d)
			w.Flush()
		}
		return err
	}
}

func handleAddDevice(cmd *cobra.Command, args []string) error {

	protocols, labels, err := getDeviceAttributes()
	if err != nil {
		return err
	}

	msg, err := getCoreMetaDataService().AddDevice(deviceName, deviceDescription, deviceService, deviceProfile,
		deviceAdminState, deviceOperState, labels, deviceLocation, protocols)
	if msg != nil {
		fmt.Println(msg)
	}
	return err
}

func handleListDevices(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().ListAllDevicesJSON(offset, limit, labels)
		if err != nil {
			return err
		}
		fmt.Print(json)
	} else {
		devices, err := getCoreMetaDataService().ListAllDevices(offset, limit, getLabels())
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			fmt.Println("No devices available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printDeviceTableHeader(w)
		for _, device := range devices {
			printDevice(w, &device)
		}
		w.Flush()
	}
	return nil
}

func getDeviceAttributes() (protocols map[string]dtos.ProtocolProperties, labels []string, err error) {

	labels = getLabels()
	if deviceProtocols != "" {
		err = jsonpkg.Unmarshal([]byte(deviceProtocols), &protocols)
		if err != nil {
			err = errors.New(`please specify protocols using a JSON value. Example: --protocols "{\"modbus-tcp\":{\"Address\": \"localhost\",\"Port\": \"1234\" }}""`)
		}
	}
	return
}

func printDeviceTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tDescription\tServiceName\tProfileName\tAdminState\tOperatingState\tLastReported\tLastConnected\tLabels\tLocation\tAutoEvents\tProtocols")
	} else {
		fmt.Fprintln(w, "Name\tDescription\tServiceName\tProfileName\tLabels\tAutoEvents")
	}

}

func printDevice(w *tabwriter.Writer, d *dtos.Device) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			d.Id,
			d.Name,
			d.Description,
			d.ServiceName,
			d.ProfileName,
			d.AdminState,
			d.OperatingState,
			time.Unix(0, d.LastReported).Format(time.RFC822),
			time.Unix(0, d.LastConnected).Format(time.RFC822),
			d.Labels,
			d.Location,
			d.AutoEvents,
			d.Protocols)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
			d.Name,
			d.Description,
			d.ServiceName,
			d.ProfileName,
			d.Labels,
			d.AutoEvents)
	}
}
