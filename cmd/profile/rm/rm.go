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
	"context"
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/urlclient"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/spf13/cobra"
)

// NewCommand return the rm profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [profile name|ID]",
		Short: "Remove profile by name or ID",
		Long:  `Removes the device profile given a device profile name or ID.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Checking for args
			if len(args) == 0 {
				fmt.Printf("Error: No profile ID/Name provided.\n")
				return
			}

			deviceID := args[0]

			ctx, _ := context.WithCancel(context.Background())

			url := config.Conf.MetadataService.Protocol + "://" +
				config.Conf.MetadataService.Host + ":" +
				config.Conf.MetadataService.Port

			mdc := metadata.NewDeviceProfileClient(
				urlclient.New(
					ctx,
					clients.CoreMetaDataServiceKey,
					clients.ApiDeviceProfileRoute,
					15000,
					url +  clients.ApiDeviceProfileRoute,
				),
			)

			//err := mdc.Delete(ctx, deviceID)

			err := mdc.DeleteByName(ctx, deviceID)

			if err == nil {
				fmt.Printf("Removed: %s\n", deviceID)
				return
			}

			err = mdc.Delete(ctx, deviceID)
			if err != nil {
				fmt.Printf("Not removed: %s\n", deviceID)
				fmt.Println(err)
				return
			}
		},
	}
	return cmd
}
