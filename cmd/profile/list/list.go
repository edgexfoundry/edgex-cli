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
	"html/template"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

const profileTemplete = "Profile ID\tProfile Name\tManufacturer\tModel\tCreated\tModified\n" +
	"{{range .}}" +
	"{{.Id}}\t{{.Name}}\t{{.Manufacturer}}\t{{.Model}}\t{{DisplayDuration .Created}}\t{{DisplayDuration .Modified}}\n" +
	"{{end}}"

// NewCommand return the list profiles command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of device profiles",
		Long:  `Returns the list of device profiles currently in the core-metadata database.`,
		RunE: listHandler,
	}
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	url := config.Conf.Clients["Metadata"].Url()
	mdc := metadata.NewDeviceProfileClient(
		local.New(url + clients.ApiDeviceProfileRoute),
	)

	profiles, err := mdc.DeviceProfiles(context.Background())
	if err != nil {
		return
	}

	// TODO: Add commands and resources? They both are lists, so we need to think about how to display them
	formatter := formatters.NewFormatter(profileTemplete, template.FuncMap{"DisplayDuration": utils.DisplayDuration})
	err = formatter.Write(profiles)
	return
}
