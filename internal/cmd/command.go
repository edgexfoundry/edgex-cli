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
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/spf13/cobra"
)

var commandDeviceName, commandName string
var pushEvent, noReturnEvent bool
var requestBody, requestFile string

func init() {
	commandCmd := initCommandCommand()
	initReadCommand(commandCmd)
	initWriteCommand(commandCmd)
	initListCommand(commandCmd)
}

func initCommandCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "command",
		Short:        "Read, write and list commands [Core Command]",
		Long:         "",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	return cmd
}

func initReadCommand(cmd *cobra.Command) {
	var readCmd = &cobra.Command{
		Use:          "read",
		Short:        "Issue a read command to the specified device",
		Long:         "Issue a read command to the specified device",
		RunE:         handleReadCommand,
		SilenceUsage: true,
	}
	readCmd.Flags().StringVarP(&commandDeviceName, "device", "d", "", "specify the name of device")
	readCmd.Flags().StringVarP(&commandName, "command", "c", "", "specify the name of the command to be executed")
	readCmd.Flags().BoolVarP(&pushEvent, "pushevent", "p", false, "if set, a successful read command will result in an event being pushed to the EdgeX system")
	readCmd.Flags().BoolVarP(&noReturnEvent, "noreturnevent", "r", false, "if set, there will be no event returned in the HTTP response")
	readCmd.MarkFlagRequired("device")
	readCmd.MarkFlagRequired("command")
	cmd.AddCommand(readCmd)
	addFormatFlags(readCmd)
}

func initWriteCommand(cmd *cobra.Command) {
	var writeCmd = &cobra.Command{
		Use:          "write",
		Short:        "Issue a write command to the specified device",
		Long:         "Issue a write command to the specified device",
		RunE:         handleWriteCommand,
		SilenceUsage: true,
	}
	writeCmd.Flags().StringVarP(&commandDeviceName, "device", "d", "", "specify the name of the device")
	writeCmd.Flags().StringVarP(&commandName, "command", "c", "", "specify the name of the command to be executed")
	writeCmd.Flags().StringVarP(&requestBody, "body", "b", "", "specify the write command's request body, which provides the value(s) being written to the device")
	writeCmd.Flags().StringVarP(&requestFile, "file", "f", "", "specify a file containing the write command's request body, which provides the value(s) being written to the device")
	writeCmd.MarkFlagRequired("device")
	writeCmd.MarkFlagRequired("command")
	cmd.AddCommand(writeCmd)
	addFormatFlags(writeCmd)
}

func initListCommand(cmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:          "list",
		Short:        "A list of device supported commands",
		Long:         "Returns a paginated list of all supported device commands, optionally filtered by device name",
		RunE:         handleListCommand,
		SilenceUsage: true,
	}
	listCmd.Flags().StringVarP(&commandDeviceName, "device", "d", "", "list commands specified by device name")
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
	addFormatFlags(listCmd)
}

func handleReadCommand(cmd *cobra.Command, args []string) error {

	// provide error message when -returnEvent and -pushEvent flag are all set to "no"
	if !pushEvent && noReturnEvent {
		fmt.Println("Nothing to do. Please remove -noreturnevent flag or set -pushevent flag.")
		return nil
	}

	// parse flags to "yes"/"no" strings
	dsPushEvent := boolToString(pushEvent)
	dsReturnEvent := boolToString(!noReturnEvent)

	response, err := getCoreCommandService().GetCommandClient().IssueGetCommandByName(context.Background(), commandDeviceName, commandName, dsPushEvent, dsReturnEvent)
	if err != nil {
		return err
	}

	// return early, if no response received
	if response == nil {
		fmt.Println("Request successful. Pushed results to EdgeX system.")
		return nil
	}

	// print READ command's output with one of these formats: JSON or table
	if json {
		stringifiedResponse, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", stringifiedResponse)
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

func handleWriteCommand(cmd *cobra.Command, args []string) error {
	// issue WRITE command's request body by one of these options: inline body or using file
	if (requestBody != "" && requestFile != "") || (requestBody == "" && requestFile == "") {
		return errors.New("please specify request data using one of the provided ways: --body or --file")
	}

	if requestFile != "" {
		content, err := ioutil.ReadFile(requestFile)
		if err != nil {
			return err
		}
		requestBody = string(content)
	}

	var settings map[string]string
	err := jsonpkg.Unmarshal([]byte(requestBody), &settings)
	if err != nil {
		return err
	}

	// issue write command

	response, err := getCoreCommandService().GetCommandClient().IssueSetCommandByName(context.Background(), commandDeviceName, commandName, settings)
	if err != nil {
		return err
	}

	// print WRITE command's output with one of these formats: JSON or string
	if json {
		stringifiedResponse, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(stringifiedResponse))
	} else {
		fmt.Printf("apiVersion: %s,statusCode: %d\n", response.ApiVersion, response.StatusCode)
	}
	return nil
}

func handleListCommand(cmd *cobra.Command, args []string) error {
	// issue list commands with specified device name
	if deviceName != "" {
		response, err := getCoreCommandService().GetCommandClient().DeviceCoreCommandsByDeviceName(context.Background(), commandDeviceName)
		if err != nil {
			return err
		}

		// print LIST commands with one of these formats: JSON or table
		if json {
			stringified, err := jsonpkg.Marshal(response)
			if err != nil {
				return err
			}
			fmt.Println(string(stringified))
		} else {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			fmt.Fprintln(w, "Name\tDevice Name\tProfile Name\tMethods\tURL")
			for _, command := range response.DeviceCoreCommand.CoreCommands {
				methods := methodsToString(command)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					command.Name, response.DeviceCoreCommand.DeviceName, response.DeviceCoreCommand.ProfileName, methods, command.Url+command.Path)
			}
			w.Flush()
		}

	} else {
		// issue list all commands, optionally specifying a limit and offset

		response, err := getCoreCommandService().GetCommandClient().AllDeviceCoreCommands(context.Background(), offset, limit)
		if err != nil {
			return err
		}

		// print LIST command's output with one of these formats: JSON or table
		if json {
			stringified, err := jsonpkg.Marshal(response)
			if err != nil {
				return err
			}

			fmt.Println(string(stringified))
		} else {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
			fmt.Fprintln(w, "Name\tDevice Name\tProfile Name\tMethods\tURL")
			for _, device := range response.DeviceCoreCommands {
				for _, command := range device.CoreCommands {
					methods := methodsToString(command)
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
						command.Name, device.DeviceName, device.ProfileName, methods, command.Url+command.Path)
				}
			}
			w.Flush()
		}
	}

	return nil
}

// used by list command when it shows in table format
func methodsToString(command dtos.CoreCommand) string {
	if command.Get && command.Set {
		return "Get, Put"
	} else if command.Get {
		return "Get"
	} else {
		return "Put"
	}
}

// used by read command when it specified dsPushEvent or dsReturnEvent
func boolToString(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}
