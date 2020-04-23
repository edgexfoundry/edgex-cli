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

package list

import (
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand returns the list device services command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists existing devices services",
		Long:  `Return the list fo current device services.`,
		RunE: func(cmd *cobra.Command, args []string) (err error){
			url := config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute
			var deviceServices []models.DeviceService
			err = client.ListHelper(url, &deviceServices)
			if err != nil {
				return
			}

			pw := viper.Get("writer").(io.WriteCloser)
			w := new(tabwriter.Writer)
			w.Init(pw, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "Service ID", "Service Name", "Created", "Operating State")
			for _, deviceService := range deviceServices {
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t\n",
					deviceService.Id,
					deviceService.Name,
					utils.DisplayDuration(deviceService.Created),
					deviceService.OperatingState,
				)
			}
			w.Flush()
			return
		},
	}
	return cmd
}
