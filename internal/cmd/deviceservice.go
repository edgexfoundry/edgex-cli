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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
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
	initGetDeviceServiceByNameCommand(cmd)
}

// initRmDeviceServiceCommand implements the DELETE ​/deviceservice/name​/{name} endpoint
// "Delete a device service by its unique name"
func initRmDeviceServiceCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:          "rm",
		Short:        "Remove a device service",
		Long:         "Removes a device service from the core-metadata database",
		RunE:         handleRmDeviceService,
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

// initAddDeviceServiceCommand implements the POST ​/deviceservice/ endpoint
// "Add a new DeviceService - name must be unique."
func initAddDeviceServiceCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "add",
		Short: "Add a new device service",
		Long: `Add a new device service

Example: 
 edgex-cli deviceservice add -n TestDeviceService -b "http://localhost:51234" -l label-one,label-two,label-three
		`,
		RunE:         handleAddDeviceService,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device name")
	add.Flags().StringVarP(&deviceServiceDescription, "description", "d", "", "Device service description")
	add.Flags().StringVarP(&deviceServiceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	add.Flags().StringVarP(&deviceServiceBaseAddress, "base-address", "b", "", "Base URL for the service")
	addLabelsFlag(add)

	add.MarkFlagRequired("name")
	add.MarkFlagRequired("base-address")
	cmd.AddCommand(add)
}

// initListDeviceServiceCommand implements the GET ​/deviceservice/all endpoint
// "Given the entire range of device services sorted by last modified descending, returns a portion of
// that range according to the offset and limit parameters. Device services may also be filtered by label."
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

// initUpdateDeviceServiceCommand implements the PATCH ​/deviceservice endpoint
// "Allows updates to an existing device service"
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
		RunE:         handleUpdateDeviceService,
		SilenceUsage: true,
	}
	updateCmd.Flags().StringVarP(&deviceServiceName, "name", "n", "", "Device service name")
	updateCmd.Flags().StringVarP(&deviceServiceID, "id", "i", "", "Device service ID")
	updateCmd.Flags().StringVarP(&deviceServiceDescription, "description", "d", "", "Device service description")
	updateCmd.Flags().StringVarP(&deviceServiceAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	updateCmd.Flags().StringVarP(&deviceServiceBaseAddress, "base-address", "b", "", "Base URL for the service")
	addLabelsFlag(updateCmd)

	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("id")
	cmd.AddCommand(updateCmd)
}

// initGetDeviceServiceByNameCommand implements the GET /deviceservice/name/{name} endpoint
// "Returns a device service by its unique name"
func initGetDeviceServiceByNameCommand(cmd *cobra.Command) {
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

func handleAddDeviceService(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceServiceClient()

	err := validateAdminState(deviceServiceAdminState)
	if err != nil {
		return err
	}

	var req = requests.NewAddDeviceServiceRequest(dtos.DeviceService{
		Name:        deviceServiceName,
		Description: deviceServiceDescription,
		Labels:      getLabels(),
		BaseAddress: deviceServiceBaseAddress,
		AdminState:  deviceServiceAdminState,
	})

	response, err := client.Add(context.Background(), []requests.AddDeviceServiceRequest{req})

	if response != nil {
		fmt.Println(response[0])
	}
	return err

}

func handleRmDeviceService(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceServiceClient()
	response, err := client.DeleteByName(context.Background(), deviceServiceName)
	if err == nil {
		fmt.Println(response)
	}
	return err
}

func handleGetDeviceServiceByName(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetDeviceServiceClient()

	response, err := client.DeviceServiceByName(context.Background(), deviceServiceName)
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
		printServiceTableHeader(w)
		printService(w, &response.Service)
		w.Flush()
	}
	return nil

}

func handleListDeviceServices(cmd *cobra.Command, args []string) error {

	client := getCoreMetaDataService().GetDeviceServiceClient()
	response, err := client.AllDeviceServices(context.Background(), getLabels(), offset, limit)
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

		if len(response.Services) == 0 {
			fmt.Println("No device services available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printServiceTableHeader(w)
		for _, p := range response.Services {
			printService(w, &p)
		}
		w.Flush()
	}
	return nil
}

func handleUpdateDeviceService(cmd *cobra.Command, args []string) (err error) {
	client := getCoreMetaDataService().GetDeviceServiceClient()

	var name, id, description, baseAddress, adminState *string
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
		baseAddress = &deviceServiceBaseAddress
	}
	if deviceServiceAdminState != "" {
		adminState = &deviceServiceAdminState
		err := validateAdminState(deviceServiceAdminState)
		if err != nil {
			return err
		}
	}

	var req = requests.NewUpdateDeviceServiceRequest(dtos.UpdateDeviceService{
		Name:        name,
		Id:          id,
		Description: description,
		Labels:      getLabels(),
		BaseAddress: baseAddress,
		AdminState:  adminState,
	})

	response, err := client.Update(context.Background(), []requests.UpdateDeviceServiceRequest{req})

	if response != nil {
		fmt.Println(response[0])
	}
	return err
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
			getRFC822Time(deviceService.LastConnected),
			getRFC822Time(deviceService.LastReported),
			getRFC822Time(deviceService.Modified))
	} else {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			deviceService.Name,
			deviceService.BaseAddress,
			deviceService.Description)
	}

}
