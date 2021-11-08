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

var intervalInterval, intervalName, intervalStart, intervalEnd, intervalId string

func init() {
	var cmd = &cobra.Command{
		Use:          "interval",
		Short:        "Add, get and list intervals [Support Scheduler]",
		Long:         "Add, get and list intervals [Support Scheduler]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initListIntervalCommand(cmd)
	initGetIntervalByNameCommand(cmd)
	initAddIntervalCommand(cmd)
	initRmIntervalCommand(cmd)
	initUpdateIntervalCommand(cmd)
}

// initListIntervalCommand implements support for the GET /interval/all endpoint
// "Given the entire range of intervals sorted by last modified descending,
// returns a portion of that range according to the offset and limit parameters."
func initListIntervalCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all intervals",
		Long:         "List all intervals",
		RunE:         handleListIntervals,
		SilenceUsage: true,
	}
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
}

// initGetIntervalByNameCommand implements support for the GET /interval/name/{name} endpoint
// "Returns an interval according to the specified name"
func initGetIntervalByNameCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Return an interval by name",
		Long:         `Return an interval by name`,
		RunE:         handleGetIntervalByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&intervalName, "name", "n", "", "Interval name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

// initAddIntervalCommand implements support for the POST /interval endpoint
// "Add one or more new Intervals - name on each request must be unique."
func initAddIntervalCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "add",
		Short:        "Add an interval",
		Long:         "Add an interval",
		Example:      `  edgex-cli interval add -n "hourly" -i "1h"`,
		RunE:         handleAddInterval,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&intervalName, "name", "n", "", "Non-database identifier for an interval (*must be unique)")
	add.Flags().StringVarP(&intervalInterval, "interval", "i", "", "Interval indicates how often the specific resource needs to be polled (e.g. 100ms, 24h)")
	add.Flags().StringVarP(&intervalStart, "start", "s", "", "Start time in ISO 8601 format YYYYMMDD'T'HHmmss")
	add.Flags().StringVarP(&intervalEnd, "end", "e", "", "End time in ISO 8601 format YYYYMMDD'T'HHmmss")

	add.MarkFlagRequired("name")
	add.MarkFlagRequired("interval")
	cmd.AddCommand(add)
}

// initRmIntervalCommand implements the DELETE /interval/name/{name}
// "Deletes an interval according to the specified name. Associated actions will also be deleted."
func initRmIntervalCommand(cmd *cobra.Command) {
	var rm = &cobra.Command{
		Use:          "rm",
		Short:        "Delete a named interval and associated interval actions",
		Long:         "Delete a named interval and associated interval actions",
		RunE:         handleRmInterval,
		SilenceUsage: true,
	}
	rm.Flags().StringVarP(&intervalName, "name", "n", "", "Interval name")
	rm.MarkFlagRequired("name")
	cmd.AddCommand(rm)
}

// initUpdateIntervalCommand implements support for the PATCH /interval endpoint
// "Update one or more existing Intervals"
func initUpdateIntervalCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "update",
		Short:        "Update an interval",
		Long:         "Update an interval, specifying either ID or name ",
		RunE:         handleUpdateInterval,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&intervalId, "id", "", "", "Uniquely identifies the interval, either id or name should be specified.")
	add.Flags().StringVarP(&intervalName, "name", "n", "", "Non-database identifier for an interval (*must be unique), either id or name should be specified")
	add.Flags().StringVarP(&intervalInterval, "interval", "i", "", "Interval indicates how often the specific resource needs to be polled (e.g. 100ms, 24h)")
	add.Flags().StringVarP(&intervalStart, "start", "s", "", "Start time in ISO 8601 format YYYYMMDD'T'HHmmss")
	add.Flags().StringVarP(&intervalEnd, "end", "e", "", "End time in ISO 8601 format YYYYMMDD'T'HHmmss")

	cmd.AddCommand(add)
}

func handleUpdateInterval(cmd *cobra.Command, args []string) error {
	var name, id, start, end, interval *string

	client := getSupportSchedulerService().GetIntervalClient()
	if intervalName != "" {
		name = &intervalName
	}
	if intervalId != "" {
		id = &intervalId
	}
	if name == nil && id == nil {
		return errors.New("either id or name should be specified")
	}
	if intervalStart != "" {
		start = &intervalStart
	}
	if intervalEnd != "" {
		end = &intervalEnd
	}
	if intervalInterval != "" {
		interval = &intervalInterval
	}
	var req = requests.NewUpdateIntervalRequest(dtos.UpdateInterval{
		Name:     name,
		Id:       id,
		Start:    start,
		End:      end,
		Interval: interval})
	response, err := client.Update(context.Background(), []requests.UpdateIntervalRequest{req})
	if response != nil {
		fmt.Println(response[0])
	}
	return err
}

func handleRmInterval(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalClient()
	response, err := client.DeleteIntervalByName(context.Background(), intervalName)
	if err == nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleAddInterval(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalClient()
	var req = requests.NewAddIntervalRequest(dtos.Interval{
		Name:     intervalName,
		Interval: intervalInterval,
		Start:    intervalStart,
		End:      intervalEnd})
	response, err := client.Add(context.Background(), []requests.AddIntervalRequest{req})
	if err != nil {
		return err
	}
	if response != nil {
		fmt.Println(response[0])
	}
	return err
}

func handleGetIntervalByName(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalClient()
	response, err := client.IntervalByName(context.Background(), intervalName)
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
		printIntervalTableHeader(w)
		printInterval(w, &response.Interval)
		w.Flush()
	}
	return nil
}

func handleListIntervals(cmd *cobra.Command, args []string) error {
	client := getSupportSchedulerService().GetIntervalClient()
	response, err := client.AllIntervals(context.Background(), offset, limit)
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

		if len(response.Intervals) == 0 {
			fmt.Println("No intervals available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printIntervalTableHeader(w)
		for _, n := range response.Intervals {
			printInterval(w, &n)
		}
		w.Flush()
	}
	return nil
}

func printIntervalTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tInterval\tStart\tEnd")
	} else {
		fmt.Fprintln(w, "Name\tInterval\tStart\tEnd")
	}

}

func printInterval(w *tabwriter.Writer, n *dtos.Interval) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
			n.Id,
			n.Name,
			n.Interval,
			n.Start,
			n.End)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n",
			n.Name,
			n.Interval,
			n.Start,
			n.End)
	}
}
