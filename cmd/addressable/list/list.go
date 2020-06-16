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
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
)

const addrTempl = "Id\tName\tProtocol\tHTTPMethod\tAddress\tPort\n" +
	"{{range .}}" +
		"{{.Id}}\t{{.Name}}\t{{.Protocol}}\t{{.HTTPMethod}}\t{{.Address}}\t{{.Port}}\n" +
	"{{end}}"

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all device services",
		Long:  `Return all device services sorted by id.`,
		RunE:  listHandler,
	}
	return cmd
}
func listHandler(cmd *cobra.Command, args []string) (err error) {
	url := config.Conf.Clients["Metadata"].Url() + clients.ApiAddressableRoute
	var addr []models.Addressable
	err = request.Get(url, &addr)
	if err != nil {
		return
	}
	formatter := formatters.NewFormatter(addrTempl, nil)
	err = formatter.Write(addr)
	return
}
