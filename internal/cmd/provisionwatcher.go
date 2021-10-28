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

var provisionWatcherName, provisionWatcherIdentifiers, provisionWatcherProfileName string
var provisionWatcherServiceName, provisionWatcherAdminState, provisionWatcherId string

func init() {
	var cmd = &cobra.Command{
		Use:          "provisionwatcher",
		Short:        "Add, remove, get, list and modify provison watchers [Core Metadata]",
		Long:         "Add, remove, get, list and modify provison watchers [Core Metadata]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)

	initAddProvisionWatcherCommand(cmd)
	initListProvisionWatcherCommand(cmd)
	initGetProvisionWatcherByNameCommand(cmd)
	initRmProvisionWatcherCommand(cmd)
	initUpdateProvisionWatcherCommand(cmd)

}

// initRmProvisionWatcherCommand implements the DELETE ​/provisionwatcher​/name​/{name} endpoint
// "Delete a provision watcher by its unique name"
func initRmProvisionWatcherCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:          "rm",
		Short:        "Remove a provision watcher",
		Long:         "Remove a provision watcher from the core-metadata database",
		RunE:         handleRmProvisionWatcher,
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&provisionWatcherName, "name", "n", "", "Provision watcher name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

// initListProvisionWatcherCommand implements the GET ​/provisionwatcher​/all endpoint:
// "Given the entire range of provision watchers sorted by last modified descending,
// returns a portion of that range according to the offset and limit parameters. Provision watchers may also be filtered by label."
func initListProvisionWatcherCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all provision watchers",
		Long:         "List all provision watchers, optionally specifying a limit, offset and/or label(s)",
		RunE:         handleListProvisionWatchers,
		SilenceUsage: true,
	}

	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	addLabelsFlag(listCmd)

}

// initUpdateProvisionWatcherCommand implements the PATCH ​/provisionwatcher endpoint
// "Allows updates to an existing provision watcher"
func initUpdateProvisionWatcherCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "update",
		Short: "Update a new provision watcher",
		Long: `Update a new provision watcher

Example: 
edgex-cli provisionwatcher add -n watcher -i "e69ec9b4-f164-4e09-8b1b-988fc545f9fb" --labels "newlabel" 
`,
		RunE:         handleUpdateProvisionWatcher,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&provisionWatcherName, "name", "n", "", "Provision watcher name")
	add.Flags().StringVarP(&provisionWatcherId, "id", "i", "", "Provision watcher ID")
	add.Flags().StringVarP(&provisionWatcherIdentifiers, "identifiers", "", "", "Set of key value pairs that identify property and value to watch for")
	add.Flags().StringVarP(&provisionWatcherProfileName, "profile", "p", "", "Name of the profile that should be applied to the devices available at the identifier addresses")
	add.Flags().StringVarP(&provisionWatcherServiceName, "service", "s", "", "Name of the device service that new devices will be associated to")
	add.Flags().StringVarP(&provisionWatcherAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	addLabelsFlag(add)
	add.MarkFlagRequired("name")
	add.MarkFlagRequired("id")
	cmd.AddCommand(add)
}

// initAddProvisionWatcherCommand implements the POST ​/provisionwatcher endpoint
// "Add a new ProvisionWatcher - name must be unique."
func initAddProvisionWatcherCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:   "add",
		Short: "Add a new provision watcher",
		Long: `Add a new provision watcher

Example: 
	edgex-cli provisionwatcher add -n watcher --identifiers "{\"address\":\"localhost\",\"port\":\"1234\"}" --labels "label,label2" -p "device-simple" -s "device-simple"
`,
		RunE:         handleAddProvisionWatcher,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&provisionWatcherName, "name", "n", "", "Provision watcher name")
	add.Flags().StringVarP(&provisionWatcherIdentifiers, "identifiers", "", "", "Set of key value pairs that identify property and value to watch for")
	add.Flags().StringVarP(&provisionWatcherProfileName, "profile", "p", "", "Name of the profile that should be applied to the devices available at the identifier addresses")
	add.Flags().StringVarP(&provisionWatcherServiceName, "service", "s", "", "Name of the device service that new devices will be associated to")
	add.Flags().StringVarP(&provisionWatcherAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")
	addLabelsFlag(add)
	add.MarkFlagRequired("name")
	add.MarkFlagRequired("profile")
	add.MarkFlagRequired("service")
	add.MarkFlagRequired("identifiers")
	cmd.AddCommand(add)
}

// initGetProvisionWatcherByNameCommand implements the GET ​/provisionwatcher/name/{name}
// "Returns a provision watcher by its unique name endpoint"
func initGetProvisionWatcherByNameCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Returns a provision watcher by name",
		Long:         `Returns a provision watcher by name`,
		RunE:         handleGetProvisionWatcherByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&provisionWatcherName, "name", "n", "", "Provision watcher name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleRmProvisionWatcher(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetProvisionWatcherClient()
	response, err := client.DeleteProvisionWatcherByName(context.Background(), provisionWatcherName)
	if err == nil {
		fmt.Println(response)
	}
	return err
}

func handleGetProvisionWatcherByName(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetProvisionWatcherClient()

	response, err := client.ProvisionWatcherByName(context.Background(), provisionWatcherName)
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
		printProvisionWatcherTableHeader(w)
		printProvisionWatcher(w, &response.ProvisionWatcher)
		w.Flush()
	}
	return nil
}

func handleListProvisionWatchers(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetProvisionWatcherClient()
	response, err := client.AllProvisionWatchers(context.Background(), getLabels(), offset, limit)
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

		if len(response.ProvisionWatchers) == 0 {
			fmt.Println("No provision watchers available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printProvisionWatcherTableHeader(w)
		for _, p := range response.ProvisionWatchers {
			printProvisionWatcher(w, &p)
		}
		w.Flush()
	}
	return nil
}

func handleUpdateProvisionWatcher(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetProvisionWatcherClient()

	var name, id, service, profile, adminState *string

	if provisionWatcherName != "" {
		name = &provisionWatcherName
	}
	if provisionWatcherId != "" {
		id = &provisionWatcherId
	}

	if provisionWatcherServiceName != "" {
		service = &provisionWatcherServiceName
	}

	if provisionWatcherProfileName != "" {
		profile = &provisionWatcherProfileName
	}

	if provisionWatcherAdminState != "" {
		adminState = &provisionWatcherAdminState
		err := validateAdminState(provisionWatcherAdminState)
		if err != nil {
			return err
		}
	}

	identifiers, labels, err := getProvisonWatcherAttributes()
	if err != nil {
		return err
	}

	var req = requests.NewUpdateProvisionWatcherRequest(dtos.UpdateProvisionWatcher{
		Name:        name,
		Id:          id,
		ServiceName: service,
		ProfileName: profile,
		AdminState:  adminState,
		Labels:      labels,
		Identifiers: identifiers,
	})

	response, err := client.Update(context.Background(), []requests.UpdateProvisionWatcherRequest{req})

	if response != nil {
		fmt.Println(response[0])
	}
	return err

}

func handleAddProvisionWatcher(cmd *cobra.Command, args []string) error {
	client := getCoreMetaDataService().GetProvisionWatcherClient()

	err := validateAdminState(provisionWatcherAdminState)
	if err != nil {
		return err
	}

	identifiers, labels, err := getProvisonWatcherAttributes()
	if err != nil {
		return err
	}

	var req = requests.NewAddProvisionWatcherRequest(dtos.ProvisionWatcher{
		Name:                provisionWatcherName,
		ServiceName:         provisionWatcherServiceName,
		ProfileName:         provisionWatcherProfileName,
		AdminState:          provisionWatcherAdminState,
		Labels:              labels,
		Identifiers:         identifiers,
		BlockingIdentifiers: nil,
	})
	response, err := client.Add(context.Background(), []requests.AddProvisionWatcherRequest{req})

	if response != nil {
		fmt.Println(response[0])
	}
	return err

}

func getProvisonWatcherAttributes() (identifiers map[string]string, labels []string, err error) {
	labels = getLabels()

	if provisionWatcherIdentifiers != "" {
		err = jsonpkg.Unmarshal([]byte(provisionWatcherIdentifiers), &identifiers)
		if err != nil {
			err = errors.New(`please specify identifiers using a JSON value. Example: --identifiers "{\"address\":\"localhost\",\"port\":\"1234\"}" `)
		}
	}
	return
}

func printProvisionWatcherTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tServiceName\tProfileName\tAdminState\tLabels\tIdentifiers\tBlockingIdentifiers\tAutoEvents")
	} else {
		fmt.Fprintln(w, "Name\tServiceName\tProfileName\tLabels\tIdentifiers")
	}

}

func printProvisionWatcher(w *tabwriter.Writer, d *dtos.ProvisionWatcher) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			d.Id,
			d.Name,
			d.ServiceName,
			d.ProfileName,
			d.AdminState,
			d.Labels,
			d.Identifiers,
			d.BlockingIdentifiers,
			d.AutoEvents)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
			d.Name,
			d.ServiceName,
			d.ProfileName,
			d.Labels,
			d.Identifiers)
	}
}
