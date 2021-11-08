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

var intervalActionName, intervalActionIntervalName, intervalActionAddress, intervalActionId string
var intervalActionContent, intervalActionContentType, intervalActionAdminState string

func init() {
	var cmd = &cobra.Command{
		Use:          "intervalaction",
		Short:        "Get, list, update and remove interval actions [Support Scheduler]",
		Long:         "Get, list, update and remove interval actions [Support Scheduler]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initAddIntervalActionCommand(cmd)
	initListIntervalActionCommand(cmd)
	initGetIntervalActionByNameCommand(cmd)
	initRmIntervalActionCommand(cmd)
	initUpdateIntervalActionCommand(cmd)
}

// initListIntervalActionCommand implements support for the GET /intervalaction/all endpoint
// "Given the entire range of interval actions sorted by last modified descending,
// returns a portion of that range according to the offset and limit parameters."
func initListIntervalActionCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all interval actions",
		Long:         "List all interval actions",
		RunE:         handleListIntervalActions,
		SilenceUsage: true,
	}
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
}

// initGetIntervalActionByNameCommand implements support for the GET /intervalaction/name/{name} endpoint
// "Returns an interval action according to the specified name"
func initGetIntervalActionByNameCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Return an interval action by name",
		Long:         `Return an interval action by name`,
		RunE:         handleGetIntervalActionByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&intervalActionName, "name", "n", "", "Interval action name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

// initAddIntervalActionCommand implements the POST /intervalaction endpoint
// "Adds one or more notifications to be sent."
func initAddIntervalActionCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "add",
		Short:        "Add an interval action",
		Long:         "Add an interval action",
		Example:      `  edgex-cli intervalaction add -n "name01" -i "midnight" -a "{\"type\": \"REST\", \"host\": \"192.168.0.102\", \"port\": 8080, \"httpMethod\": \"GET\"}"`,
		RunE:         handleAddIntervalAction,
		SilenceUsage: true,
	}

	add.Flags().StringVarP(&intervalActionName, "name", "n", "", "Interval action name")
	add.Flags().StringVarP(&intervalActionIntervalName, "interval", "i", "", "Name of the interval associated with this action")
	add.Flags().StringVarP(&intervalActionAddress, "address", "a", "", "JSON representation of the address information")
	add.Flags().StringVarP(&intervalActionContent, "content", "c", "", "Interval action content")
	add.Flags().StringVarP(&intervalActionContentType, "content-type", "t", "", "Interval action content type  (i.e. text/html, application/json)")
	add.Flags().StringVarP(&intervalActionAdminState, "admin-state", "", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")

	addLabelsFlag(add)
	add.MarkFlagRequired("name")
	add.MarkFlagRequired("interval")
	add.MarkFlagRequired("address")
	cmd.AddCommand(add)
}

// initRmIntervalActionCommand implements the DELETE /intervalaction/name/{name}
// "Deletes an interval action by name"
func initRmIntervalActionCommand(cmd *cobra.Command) {
	var rm = &cobra.Command{
		Use:          "rm",
		Short:        "Delete an interval action by name",
		Long:         "Delete an interval action by name",
		RunE:         handleRmIntervalAction,
		SilenceUsage: true,
	}
	rm.Flags().StringVarP(&intervalActionName, "name", "n", "", "Interval action name")
	rm.MarkFlagRequired("name")
	cmd.AddCommand(rm)
}

// initUpdateIntervalActionCommand implements support for the PATCH /intervalaction endpoint
// "Update one or more existing IntervalActions"
func initUpdateIntervalActionCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "update",
		Short:        "Update an interval action",
		Long:         "Update an interval action, specifying either ID or name ",
		RunE:         handleUpdateIntervalAction,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&intervalActionId, "id", "", "", "Uniquely identifies the interval action, either id or name should be specified.")
	add.Flags().StringVarP(&intervalActionName, "name", "n", "", "Interval action name")
	add.Flags().StringVarP(&intervalActionIntervalName, "interval", "i", "", "Name of the interval associated with this action")
	add.Flags().StringVarP(&intervalActionAddress, "address", "a", "", "JSON representation of the address information")
	add.Flags().StringVarP(&intervalActionContent, "content", "c", "", "Interval action content")
	add.Flags().StringVarP(&intervalActionContentType, "content-type", "t", "", "Interval action content type  (i.e. text/html, application/json)")
	add.Flags().StringVarP(&intervalActionAdminState, "admin-state", "", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")

	cmd.AddCommand(add)
}

func handleUpdateIntervalAction(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalActionClient()

	var name, id, intervalName, content, contentType, adminState *string
	var address *dtos.Address

	if intervalActionId != "" {
		id = &intervalActionId
	}
	if intervalActionName != "" {
		name = &intervalActionName
	}
	if name == nil && id == nil {
		return errors.New("either id or name should be specified")
	}
	if intervalActionIntervalName != "" {
		intervalName = &intervalActionIntervalName
	}
	if intervalActionAddress != "" {
		address = new(dtos.Address)
		err := jsonpkg.Unmarshal([]byte(intervalActionAddress), address)
		if err != nil {
			err = fmt.Errorf("channels JSON object array invalid (%v)", err)
		}
		return err
	}
	if intervalActionContent != "" {
		content = &intervalActionContent
	}
	if intervalActionContentType != "" {
		contentType = &intervalActionContentType
	}
	if intervalActionAdminState != "" {
		err := validateAdminState(intervalActionAdminState)
		if err != nil {
			return err
		}
		adminState = &intervalActionAdminState
	}

	var req = requests.NewUpdateIntervalActionRequest(dtos.UpdateIntervalAction{
		Name:         name,
		Id:           id,
		IntervalName: intervalName,
		Content:      content,
		ContentType:  contentType,
		Address:      address,
		AdminState:   adminState})

	response, err := client.Update(context.Background(), []requests.UpdateIntervalActionRequest{req})
	if response != nil {
		fmt.Println(response[0])
	}
	return err

}

func handleAddIntervalAction(cmd *cobra.Command, args []string) error {
	var address dtos.Address

	client := getSupportSchedulerService().GetIntervalActionClient()
	err := validateAdminState(intervalActionAdminState)
	if err != nil {
		return err
	}
	if intervalActionAddress != "" {
		err = jsonpkg.Unmarshal([]byte(intervalActionAddress), &address)
		if err != nil {
			err = fmt.Errorf("channels JSON object array invalid (%v)", err)
			return err
		}
	}

	var req = requests.NewAddIntervalActionRequest(dtos.IntervalAction{
		Name:         intervalActionName,
		IntervalName: intervalActionIntervalName,
		Address:      address,
		Content:      intervalActionContent,
		ContentType:  intervalActionContentType,

		AdminState: intervalActionAdminState})
	response, err := client.Add(context.Background(), []requests.AddIntervalActionRequest{req})

	if err == nil && response != nil {
		fmt.Println(response[0])
	}

	return err
}

func handleRmIntervalAction(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalActionClient()

	response, err := client.DeleteIntervalActionByName(context.Background(), intervalActionName)
	if err == nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleGetIntervalActionByName(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalActionClient()
	response, err := client.IntervalActionByName(context.Background(), intervalActionName)
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
		printIntervalActionTableHeader(w)
		printIntervalAction(w, &response.Action)
		w.Flush()
	}
	return nil
}

func handleListIntervalActions(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalActionClient()
	response, err := client.AllIntervalActions(context.Background(), offset, limit)
	if err != nil {
		return err
	}

	if json {
		result, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}
		fmt.Print(string(result))
	} else {

		if len(response.Actions) == 0 {
			fmt.Println("No interval actions available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printIntervalActionTableHeader(w)
		for _, n := range response.Actions {
			printIntervalAction(w, &n)
		}
		w.Flush()
	}
	return nil
}

func printIntervalActionTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tInterval\tAddress\tContent\tContentType\tAdminState\tCreated\tUpdated")
	} else {
		fmt.Fprintln(w, "Name\tInterval\tAddress\tContent\tContentType")
	}
}

func printIntervalAction(w *tabwriter.Writer, n *dtos.IntervalAction) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			n.Id,
			n.Name,
			n.IntervalName,
			n.Address,
			n.Content,
			n.ContentType,
			n.AdminState,
			getRFC822Time(n.Created),
			getRFC822Time(n.Modified))
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
			n.Name,
			n.IntervalName,
			n.Address,
			n.Content,
			n.ContentType)
	}
}
