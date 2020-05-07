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

package list

import (
	"context"
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"


	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand return the list profiles command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of device profiles",
		Long:  `Returns the list of device profiles currently in the core-metadata database.`,
		RunE: func(cmd *cobra.Command, args []string) (err error){

			url := config.Conf.Clients["Metadata"].Url()
			mdc := metadata.NewDeviceProfileClient(
				local.New(url+clients.ApiDeviceProfileRoute),
			)

			profiles, err := mdc.DeviceProfiles(context.Background())
			if err != nil {
				return
			}

			// TODO: Add commands and resources? They both are lists, so we need to think about how to display them

			pw := viper.Get("writer").(io.WriteCloser)
			w := new(tabwriter.Writer)
			w.Init(pw, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "Profile ID", "Profile Name", "Created", "Modified", "Manufacturer", "Model")
			for _, device := range profiles {
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%v\t\n",
					device.Id,
					device.Name,
					utils.DisplayDuration(device.Created),
					utils.DisplayDuration(device.Modified),
					device.Manufacturer,
					device.Model,
				)
			}
			w.Flush()
			return
		},
	}
	return cmd
}
