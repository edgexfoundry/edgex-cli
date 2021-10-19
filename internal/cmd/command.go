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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var deviceName, commandName string
var pushEvent, noReturnEvent bool

func init() {
	commandCmd := initCommandCommand()
	initReadCommand(commandCmd)
}

func initCommandCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "command",
		Short:        "read, write and list commands",
		Long:         ``,
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	addFormatFlags(cmd)
	return cmd
}

func initReadCommand(cmd *cobra.Command) {
	var readCmd = &cobra.Command{
		Use:          "read",
		Short:        "read command referenced by the command name and device name",
		Long:         `Issue the specified read command referenced by the command name to the device/sensor that is also referenced by name`,
		RunE:         handleReadCommand,
		SilenceUsage: true,
	}
	readCmd.Flags().StringVarP(&deviceName, "device", "d", "", "Specify the name of device")
	readCmd.Flags().StringVarP(&commandName, "command", "c", "", "Specify the name of the command to be executed")
	readCmd.Flags().BoolVarP(&pushEvent, "pushevent", "p", false, "If set, a successful GET will result in an event being pushed to the EdgeX system")
	readCmd.Flags().BoolVarP(&noReturnEvent, "noreturnevent", "r", false, "If set, there will be no Event returned in the http response")
	readCmd.MarkFlagRequired("device")
	readCmd.MarkFlagRequired("command")
	cmd.AddCommand(readCmd)
	addFormatFlags(readCmd)
}

func handleReadCommand(cmd *cobra.Command, args []string) error {
	dsPushEvent := boolToString(pushEvent)
	dsReturnEvent := boolToString(!noReturnEvent)

	//issue READ command and handle output if nothing can be displayed
	response, err := getCoreCommandService().IssueReadCommand(deviceName, commandName, dsPushEvent, dsReturnEvent)
	if err != nil {
		return err
	}

	if response == nil {
		fmt.Println("Read command issued. Nothing to display. Please retry without flag -r.")
		return nil
	}

	//print READ command's output with one of these formats: JSON, verbose or table
	if json || verbose {
		stringifiedResponse, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		if verbose {
			url := getCoreCommandService().GetReadEndpoint(deviceName, commandName, dsPushEvent, dsReturnEvent)
			fmt.Printf("Result:%s\nURL: %s\n", string(stringifiedResponse), url)
		} else {
			fmt.Printf("%s\n", stringifiedResponse)
		}
	} else {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		fmt.Fprintln(w, "Command Name\tDevice Name\tProfile Name\tValue Type\tValue")
		for _, reading := range response.Event.Readings {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				commandName, reading.DeviceName, reading.ProfileName, reading.ValueType, reading.Value)
		}
		w.Flush()
	}
	return nil
}

//using by READ when it specified dsPushEvent or dsReturnEvent
func boolToString(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}
