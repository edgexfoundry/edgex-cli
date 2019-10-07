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

// NewCommand return the rm profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [profile name|ID]",
		Short: "Remove profile by name or ID",
		Long:  `Removes the device profile given a device profile name or ID.`,
		Run: func(cmd *cobra.Command, args []string) {

			verbose, _ := cmd.Flags().GetBool("verbose")

			deviceID := args[0]
			respBody, err := client.DeleteItem(deviceID, config.Conf.MetadataService.DeviceProfileByIDRoute, config.Conf.MetadataService.DeviceProfileBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			// Display Results
			if string(respBody) == "true" {
				fmt.Printf("Removed: %s\n", deviceID)
			} else {
				fmt.Printf("Remove Unsuccessful!\n")
			}
		},
	}
	return cmd
}
