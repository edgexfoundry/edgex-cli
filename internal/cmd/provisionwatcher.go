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

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
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
	initNameProvisionWatcherCommand(cmd)
	initRmProvisionWatcherCommand(cmd)
	initUpdateProvisionWatcherCommand(cmd)

}

// initRmProvisionWatcherCommand implements the DELETE ​/provisionWatcher​/name​/{name} endpoint
// "Delete a ProvisionWatcher by name"
func initRmProvisionWatcherCommand(cmd *cobra.Command) {
	var rmcmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a provision watcher",
		Long:  "Remove a provision watcher from the core-metadata database",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := getCoreMetaDataService().RemoveProvisionWatcher(provisionWatcherName)
			if err == nil && response != nil {
				fmt.Println(response)
			}
			return err
		},
		SilenceUsage: true,
	}
	rmcmd.Flags().StringVarP(&provisionWatcherName, "name", "n", "", "Provision watcher name")
	rmcmd.MarkFlagRequired("name")
	cmd.AddCommand(rmcmd)
}

// initListProvisionWatcherCommand implements the GET ​/provisionWatcher​/all endpoint:
// "Given the entire range of ProvisionWatchers sorted by last modified descending,
// returns a portion of that range according to the offset and limit parameters.
// ProvisionWatchers may also be filtered by label."
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

// initUpdateProvisionWatcherCommand implements the PATCH ​/ProvisionWatcher endpoint
// "Allows updates to an existing ProvisionWatcher"
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

// initAddProvisionWatcherCommand implements the POST ​/ProvisionWatcher endpoint
// "Allows provisioning of a new ProvisionWatcher"
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

func initNameProvisionWatcherCommand(cmd *cobra.Command) {
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

func handleGetProvisionWatcherByName(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().GetProvisionWatcherByNameJSON(provisionWatcherName)
		if err == nil {
			fmt.Print(json)
		}
		return err
	} else {
		d, err := getCoreMetaDataService().GetProvisionWatcherByName(provisionWatcherName)
		if d != nil {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			printProvisionWatcherTableHeader(w)
			printProvisionWatcher(w, d)
			w.Flush()
		}
		return err
	}
}

func handleUpdateProvisionWatcher(cmd *cobra.Command, args []string) error {

	var name, id, service, profile, admin *string

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
		admin = &provisionWatcherAdminState
	}

	identifiers, labels, err := getProvisonWatcherAttributes()
	if err != nil {
		return err
	}

	response, err := getCoreMetaDataService().UpdateProvisionWatcher(name, id, service, profile, admin, labels, identifiers)
	if response != nil {
		fmt.Println(response)
	}
	return err

}

func handleAddProvisionWatcher(cmd *cobra.Command, args []string) error {

	identifiers, labels, err := getProvisonWatcherAttributes()
	if err != nil {
		return err
	}

	msg, err := getCoreMetaDataService().AddProvisionWatcher(
		provisionWatcherName, provisionWatcherServiceName, provisionWatcherProfileName,
		provisionWatcherAdminState, labels, identifiers)

	if msg != nil {
		fmt.Println(msg)
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

func handleListProvisionWatchers(cmd *cobra.Command, args []string) error {
	if json {
		json, _, err := getCoreMetaDataService().ListAllProvisionWatchersJSON(offset, limit, labels)
		if err != nil {
			return err
		}
		fmt.Print(json)
	} else {
		ProvisionWatchers, err := getCoreMetaDataService().ListAllProvisionWatchers(offset, limit, getLabels())
		if err != nil {
			return err
		}
		if len(ProvisionWatchers) == 0 {
			fmt.Println("No provision watchers available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printProvisionWatcherTableHeader(w)
		for _, ProvisionWatcher := range ProvisionWatchers {
			printProvisionWatcher(w, &ProvisionWatcher)
		}
		w.Flush()
	}
	return nil
}
