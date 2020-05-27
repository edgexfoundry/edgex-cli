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
	"errors"
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

var name string

// NewCommand return the rm profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [<id> | --name <profile-name>]",
		Short: "Remove profile by name or ID",
		Long:  `Removes the device profile given a device profile name or ID.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 && name == "" {
				return errors.New("no profile id/name provided.\n")
			}

			ctx := context.Background()
			mdc := metadata.NewDeviceProfileClient(
				local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute),
			)
			if name != "" {
				err = mdc.DeleteByName(ctx, name)
				if err == nil {
					fmt.Printf("Removed: %s\n", name)
				}
				return
			}

			deviceID := args[0]
			err = mdc.Delete(ctx, deviceID)
			if err == nil {
				fmt.Printf("Removed: %s\n", deviceID)
				return
			}
			return
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Delete device profile by name")
	return cmd
}
