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
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/spf13/cobra"
)

var eventLimit, eventOffset int
var eventDevice, eventProfile, eventSource, readingsValueType string
var eventAge int
var numberOfReadings int

func init() {
	eventCmd := initEventCommand()
	initListEventCommand(eventCmd)
	initCountEventCommand(eventCmd)
	initRmEventCommand(eventCmd)
	initAddEventCommand(eventCmd)
}

func initEventCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "event",
		Short:        "Add, remove and list events",
		Long:         ``,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(cmd)
	return cmd
}

func initListEventCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List events",
		Long:         `List all events, optionally specifying a limit and offset`,
		RunE:         handleListEvents,
		SilenceUsage: true,
	}

	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	listCmd.Flags().IntVarP(&eventLimit, "limit", "l", 50, "The number of items to return. Specifying -1 will return all remaining items")
	listCmd.Flags().IntVarP(&eventOffset, "offset", "o", 0, "The number of items to skip")
}

func initCountEventCommand(cmd *cobra.Command) {
	var countCmd = &cobra.Command{
		Use:          "count",
		Short:        "Count available events",
		Long:         `Count the number of events in core data, optionally filtering by device name`,
		RunE:         handleCountEvents,
		SilenceUsage: true,
	}

	countCmd.Flags().StringVarP(&eventDevice, "device", "d", "", "Device name")
	cmd.AddCommand(countCmd)
	addFormatFlags(countCmd)
}

func initRmEventCommand(cmd *cobra.Command) {
	var rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove events",
		Long: `Remove events, specifying either device name or maximum event age in milliseconds
 
'edgex-cli event rm --device {devicename}' removes all events for the specified device
'edgex-cli event rm --age {ms}' removes all events generated in the last {ms} milliseconds`,
		RunE:         handleRmEvents,
		SilenceUsage: true,
	}

	rmCmd.Flags().StringVarP(&eventDevice, "device", "d", "", "Device name")
	rmCmd.Flags().IntVarP(&eventAge, "age", "a", 0, "Event age (in milliseconds)")
	cmd.AddCommand(rmCmd)
}

func initAddEventCommand(cmd *cobra.Command) {
	var addCmd = &cobra.Command{
		Use:          "add",
		Short:        "Create an event",
		Long:         `Create an event with a specified number of random readings`,
		RunE:         handleAddEvents,
		SilenceUsage: true,
	}
	addCmd.Flags().StringVarP(&eventDevice, "device", "d", "", "Device name")
	addCmd.Flags().StringVarP(&eventProfile, "profile", "p", "", "Profile name")
	addCmd.Flags().StringVarP(&readingsValueType, "type", "t", "string", "Readings value type  [bool | string | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 ]")
	addCmd.Flags().StringVarP(&eventSource, "source", "s", "", "Event source name (ResourceName or CommandName)")
	addCmd.Flags().IntVarP(&numberOfReadings, "readings", "r", 1, "Number of sample readings to create")
	addCmd.MarkFlagRequired("device")
	addCmd.MarkFlagRequired("profile")
	addCmd.MarkFlagRequired("source")
	cmd.AddCommand(addCmd)
}

func handleAddEvents(cmd *cobra.Command, args []string) error {
	client := getCoreDataService().GetEventClient()

	if numberOfReadings < 1 {
		return errors.New("the number of readings must be at least 1")
	}

	event := dtos.NewEvent(eventProfile, eventDevice, eventSource)
	readingsValueType = strings.Title(strings.ToLower(readingsValueType))

	event.Readings = make([]dtos.BaseReading, numberOfReadings)
	for i := 0; i < numberOfReadings; i++ {
		var err error
		r64 := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
		var value interface{}
		switch readingsValueType {

		case common.ValueTypeBool:
			value = (r64&1 == 0)
		case common.ValueTypeString:
			value = "Reading " + strconv.Itoa(i)
		case common.ValueTypeUint8:
			value = uint8(r64)
		case common.ValueTypeUint16:
			value = uint16(r64)
		case common.ValueTypeUint32:
			value = uint32(r64)
		case common.ValueTypeUint64:
			value = r64
		case common.ValueTypeInt8:
			value = int8(r64)
		case common.ValueTypeInt16:
			value = int16(r64)
		case common.ValueTypeInt32:
			value = int32(r64)
		case common.ValueTypeInt64:
			value = int64(r64)
		case common.ValueTypeFloat32:
			value = float32(r64) / 100
		case common.ValueTypeFloat64:
			value = float64(r64) / 100
		default:
			return errors.New("type must be one of [bool | string | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 ]")
		}
		event.Readings[i], err = dtos.NewSimpleReading(eventProfile, eventDevice, eventSource, readingsValueType, value)
		if err != nil {
			return err
		}
	}

	response, err := client.Add(context.Background(), requests.NewAddEventRequest(event))
	if err == nil {
		fmt.Printf("Added event %v\n", response.Id)
	}
	return err
}

func handleRmEvents(cmd *cobra.Command, args []string) error {
	client := getCoreDataService().GetEventClient()

	if eventDevice != "" && eventAge != 0 {
		return errors.New("either specify device name or event age, but not both")
	} else if eventDevice != "" {
		client.DeleteByDeviceName(context.Background(), eventDevice)
	} else if eventAge != 0 {
		client.DeleteByAge(context.Background(), eventAge)
	} else {
		return errors.New("event ID, device name or event age must be specified")
	}

	return nil
}

func handleCountEvents(cmd *cobra.Command, args []string) error {
	client := getCoreDataService().GetEventClient()

	var response dtosCommon.CountResponse
	var err error

	if eventDevice != "" {
		response, err = client.EventCountByDeviceName(context.Background(), eventDevice)
	} else {
		response, err = client.EventCount(context.Background())
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
			fmt.Printf("Total %s events: %v\n", eventDevice, response.Count)
		} else {
			fmt.Printf("Total events: %v\n", response.Count)
		}
	}
	return nil
}

func handleListEvents(cmd *cobra.Command, args []string) error {
	var err error

	client := getCoreDataService().GetEventClient()
	response, err := client.AllEvents(context.Background(), eventOffset, eventLimit)
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

		if response.Events == nil || len(response.Events) == 0 {
			fmt.Println("No events available")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		if verbose {
			fmt.Fprintln(w, "Origin\tDevice\tProfile\tSource\tId\tVersionable\tReadings")
			for _, event := range response.Events {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
					time.Unix(0, event.Origin).Format(time.RFC822),
					event.DeviceName,
					event.ProfileName,
					event.SourceName,
					event.Id,
					event.Versionable,
					event.Readings)
			}

		} else {
			fmt.Fprintln(w, "Origin\tDevice\tProfile\tSource\tNumber of readings")
			for _, event := range response.Events {
				tm := time.Unix(0, event.Origin)
				sTime := tm.Format(time.RFC822)
				nReadings := 0
				if event.Readings != nil {
					nReadings = len(event.Readings)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%d\n",
					sTime, event.DeviceName, event.ProfileName, event.SourceName, nReadings)
			}
		}
		w.Flush()
	}
	return nil
}
