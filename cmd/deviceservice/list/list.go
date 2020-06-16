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
	"html/template"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

const dsTemplete = "Service ID\tService Name\tOperating State\tAdmin State\tAddressable Name\tCreated\n" +
	"{{range .}}" +
	"{{.Id}}\t{{.Name}}\t{{.OperatingState}}\t{{.AdminState}}\t{{.Addressable.Name}}\t{{DisplayDuration .Created}}\n" +
	"{{end}}"

// NewCommand returns the list device services command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists existing devices services",
		Long:  `Return the list fo current device services.`,
		RunE: listHandler,
	}
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	url := config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute
	var deviceServices []models.DeviceService
	err = request.Get(url, &deviceServices)
	if err != nil {
		return
	}

	formatter := formatters.NewFormatter(dsTemplete, template.FuncMap{"DisplayDuration": utils.DisplayDuration})
	err = formatter.Write(deviceServices)
	return
}
