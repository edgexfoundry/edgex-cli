// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rm

import (
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

// NewCommand returns the rm device service command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [device service name|ID]",
		Short: "Removes device service by name or ID",
		Long:  `Removes a device service from the core-metadata DB.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Checking for args
			if len(args) == 0 {
				fmt.Printf("Error: No device service ID/Name provided.\n")
				return
			}

			dsId := args[0]
			url:=config.Conf.Clients["Metadata"].Url()+ clients.ApiDeviceServiceRoute
			err := request.DeleteItemByIdOrName(dsId,
				config.PathId,
			    config.PathName,
			    url)

			if err != nil {
				fmt.Printf("Failed to remove Device service `%s`.\nError: %s\n", dsId, err )
				return
			}
			fmt.Printf("Removed: %s\n", dsId)
		},
	}
	return cmd
}
