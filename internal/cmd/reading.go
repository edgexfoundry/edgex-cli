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
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/spf13/cobra"
)

var readingDevice string
var readingLimit, readingOffset int

func init() {
	readingCmd := initReadingCommand()
	initListReadingCommand(readingCmd)
	initCountReadingCommand(readingCmd)
}

func initReadingCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "reading",
		Short:        "Count and list readings",
		Long:         ``,
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	return cmd
}

func initListReadingCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all readings",
		Long:         `List all readings, optionally specifying a limit and offset`,
		RunE:         handleListReadings,
		SilenceUsage: true,
	}
	listCmd.Flags().IntVarP(&readingLimit, "limit", "l", 50, "The number of items to return. Specifying -1 will return all remaining items")
	listCmd.Flags().IntVarP(&readingOffset, "offset", "o", 0, "The number of items to skip")
	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)

}

func initCountReadingCommand(cmd *cobra.Command) {
	var countCmd = &cobra.Command{
		Use:          "count",
		Short:        "Count available readings",
		Long:         `Count the number of readings in core data, optionally filtering by device name`,
		RunE:         handleCountReadings,
		SilenceUsage: true,
	}

	countCmd.Flags().StringVarP(&readingDevice, "device", "d", "", "Device name")
	cmd.AddCommand(countCmd)
	addFormatFlags(countCmd)
}

func handleCountReadings(cmd *cobra.Command, args []string) error {
	client := getCoreDataService().GetReadingClient()

	var response common.CountResponse
	var err error

	if readingDevice != "" {
		response, err = client.ReadingCountByDeviceName(context.Background(), readingDevice)
	} else {
		response, err = client.ReadingCount(context.Background())
	}

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
		if eventDevice != "" {
			fmt.Printf("Total %s readings: %v\n", readingDevice, response.Count)
		} else {
			fmt.Printf("Total readings: %v\n", response.Count)
		}
	}
	return nil
}

func handleListReadings(cmd *cobra.Command, args []string) error {
	var err error

	client := getCoreDataService().GetReadingClient()
	response, err := client.AllReadings(context.Background(), offset, limit)
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

		if response.Readings == nil || len(response.Readings) == 0 {
			fmt.Println("No readings available")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		if verbose {
			fmt.Fprintln(w, "Origin\tDeviceName\tProfileName\tValue\tValueType\tId\tMediaType\tBinaryValue")
			for _, reading := range response.Readings {

				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
					time.Unix(0, reading.Origin).Format(time.RFC822),
					reading.DeviceName,
					reading.ProfileName,
					reading.Value,
					reading.ValueType,
					reading.Id,
					reading.MediaType,
					reading.BinaryValue)
			}
		} else {
			fmt.Fprintln(w, "Origin\tDevice\tProfileName\tValue\tValueType")
			for _, reading := range response.Readings {
				tm := time.Unix(0, reading.Origin)
				sTime := tm.Format(time.RFC822)
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\n",
					sTime, reading.DeviceName, reading.ProfileName, reading.Value, reading.ValueType)

			}
		}
		w.Flush()
	}
	return nil
}
