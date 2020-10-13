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

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

var name string

// NewCommand returns the rm command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [<id> | --name <device-name>]",
		Short: "Removes device by name or ID",
		Long: `Removes a device given its name or ID. 
You can use: '$ edgex device list' to find a device's name and ID.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 && name == "" {
				return errors.New("no device id/name provided")
			}
			mdc := metadata.NewDeviceClient(
				local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute),
			)

			ctx := context.Background()
			if name != "" {
				err = mdc.DeleteByName(ctx, name)
				if err == nil {
					fmt.Printf("Removed: %s\n", name)
				}
				return err
			}

			if len(args[0]) != 0 {
				return request.DeleteByIds(mdc, args)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Delete Device by name")
	return cmd
}
