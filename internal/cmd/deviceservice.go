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
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/spf13/cobra"
)

var deviceServiceName, deviceServiceID, deviceServiceDescription, deviceServiceBaseAddress, deviceServiceAdminState string

func init() {
	var cmd = &cobra.Command{
		Use:          "deviceservice",
		Short:        "Add, remove, get, list and modify device services [Core Metadata]",
		Long:         "Add, remove, get, list and modify device services [Core Metadata]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	addFormatFlags(cmd)

	initRmDeviceServiceCommand(cmd)
	initListDeviceServiceCommand(cmd)
	initAddDeviceServiceCommand(cmd)
	initUpdateDeviceServiceCommand(cmd)
	initNameDeviceServiceCommand(cmd)
}

func initRmDeviceServiceCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a device service",
		Long:  "Removes a device service from the core-metadata database",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			response, err := getCoreMetaDataService().RemoveDeviceService(deviceServiceName)
			if response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

func initAddDeviceServiceCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "add",
		Short: "Add a new device service",
		Long: `Add a new device service - name must be unique

Example: 
 edgex-cli deviceservice add -n TestDeviceService -b "http://localhost:51234" -l label-one,label-two,label-three
		
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			response, err := getCoreMetaDataService().AddDeviceService(deviceServiceName, deviceServiceDescription, deviceServiceBaseAddress, deviceServiceAdminState, getLabels())
			if response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device name")
	add.Flags().StringVarP(&deviceServiceDescription, "description", "d", "", "Device description")
	add.Flags().StringVarP(&deviceServiceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	add.Flags().StringVarP(&deviceServiceBaseAddress, "base-address", "b", "", "Associated device profile")
	addLabelsFlag(add)

	add.MarkFlagRequired("name")
	add.MarkFlagRequired("base-address")
	cmd.AddCommand(add)
}

func initListDeviceServiceCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List device services",
		Long:         `List all device services, optionally specifying a limit, offset and/or label(s)`,
		RunE:         handleListDeviceServices,
		SilenceUsage: true,
	}

	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	addLabelsFlag(listCmd)

}

func initUpdateDeviceServiceCommand(cmd *cobra.Command) {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a new device service",
		Long: `Update an existing device service definition. 
'id' and 'deviceServiceName' must be populated in order to identify the service. 
Any other property that is populated in the request will be updated. Empty/blank properties will not be considered.

Example: 
 edgex-cli deviceservice update -n TestDeviceService -b "http://localhost:51234" -l label-one,label-two,label-three
		
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var name, id, description, address, admin *string
			if deviceServiceName != "" {
				name = &deviceServiceName
			}
			if deviceServiceID != "" {
				id = &deviceServiceID
			}
			if deviceServiceDescription != "" {
				description = &deviceServiceDescription
			}
			if deviceServiceBaseAddress != "" {
				address = &deviceServiceBaseAddress
			}
			if deviceServiceAdminState != "" {
				admin = &deviceServiceAdminState
			}

			response, err := getCoreMetaDataService().UpdateDeviceService(name, id, description, address, admin, getLabels())
			if response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	updateCmd.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device service name")
	updateCmd.Flags().StringVarP(&deviceServiceID, "id", "i", "", "Device service ID")
	updateCmd.Flags().StringVarP(&deviceServiceDescription, "description", "d", "", "Device description")
	updateCmd.Flags().StringVarP(&deviceServiceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	updateCmd.Flags().StringVarP(&deviceServiceBaseAddress, "base-address", "b", "", "Associated device profile")
	addLabelsFlag(updateCmd)

	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("id")
	cmd.AddCommand(updateCmd)
}

func initNameDeviceServiceCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Returns a device service by its unique name",
		Long:         `Returns a device service by its unique name`,
		RunE:         handleGetDeviceServiceByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device name")
	nameCmd.MarkFlagRequired("device")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleGetDeviceServiceByName(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().GetDeviceServiceByNameJSON(deviceServiceName)
		if err == nil {
			fmt.Print(json)
		}
		return err
	} else {
		service, err := getCoreMetaDataService().GetDeviceServiceByName(deviceServiceName)
		if service != nil {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			printServiceTableHeader(w)
			printService(w, service)
			w.Flush()
		}
		return err
	}
}

func printServiceTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Name\tBaseAddress\tDescription\tAdminState\tId\tLabels\tLastConnected\tLastReported\tModified")
	} else {
		fmt.Fprintln(w, "Name\tBaseAddress\tDescription")
	}

}

func printService(w *tabwriter.Writer, deviceService *dtos.DeviceService) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%s\t%s\t%s\t%s\n",
			deviceService.Name,
			deviceService.BaseAddress,
			deviceService.Description,
			deviceService.AdminState,
			deviceService.Id,
			deviceService.Labels,
			time.Unix(0, deviceService.LastConnected).Format(time.RFC822),
			time.Unix(0, deviceService.LastReported).Format(time.RFC822),
			time.Unix(0, deviceService.Modified).Format(time.RFC822))
	} else {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			deviceService.Name,
			deviceService.BaseAddress,
			deviceService.Description)
	}

}

func handleListDeviceServices(cmd *cobra.Command, args []string) error {

	if json {

		json, _, err := getCoreMetaDataService().ListAllDeviceServicesJSON(offset, limit, labels)

		if err != nil {
			return err
		}

		fmt.Print(json)

	} else {
		deviceServices, err := getCoreMetaDataService().ListAllDeviceServices(offset, limit, getLabels())

		if err != nil {
			return err
		}
		if len(deviceServices) == 0 {
			fmt.Println("No device services available")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printServiceTableHeader(w)
		for _, deviceService := range deviceServices {
			printService(w, &deviceService)
		}
		w.Flush()
	}
	return nil
}
