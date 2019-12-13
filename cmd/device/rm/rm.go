// Copyright © 2019 VMware, INC
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
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/spf13/cobra"
)

// NewCommand returns the rm command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [device name|ID]",
		Short: "Removes device by name or ID",
		Long: `Removes a device given its name or ID. 
You can use: '$ edgex device list' to find a device's name and ID.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Checking for args
			if len(args) == 0 {
				fmt.Printf("Error: No device ID/Name provided.\n")
				return
			}

			deviceID := args[0]
			respBody, err := client.DeleteItem(deviceID,
				config.Conf.MetadataService.DeviceByIDRoute,
				config.Conf.MetadataService.DeviceBySlugNameRoute,
				config.Conf.MetadataService.Port)

			if err != nil {
				fmt.Println(err)
				return
			}

			// Display Results
			if string(respBody) == "true" {
				fmt.Printf("Removed: %s\n", deviceID)
			} else {
				fmt.Printf("Error: Remove Unsuccessful: %s\n", string(respBody))
			}
		},
	}
	return cmd
}
