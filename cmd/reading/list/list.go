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
	"html/template"
	"strconv"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var limit int32

const readingTemplate = "Reading ID\tName\tDevice\tOrigin\tValue\tCreated\tModified\tPushed\n" +
	"{{range .}}" +
	"{{.Id}}\t{{.Name}}\t{{.Device}}\t{{.Origin}}\t{{.Value}}\t{{DisplayDuration .Created}}\t{{DisplayDuration .Modified}}\t{{DisplayDuration .Pushed}}\n" +
	"{{end}}"

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all device readings",
		Long:  `Return all device readings.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: listHandler,
	}
	cmd.Flags().Int32VarP(&limit, "limit", "l", 0, "Limit number of results")
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	var url string
	if len(args) > 0 {
		var limitUrl string
		device := args[0]
		if limit > 0 {
			limitUrl = strconv.FormatInt(int64(limit), 10)
		} else {
			limitUrl = strconv.FormatInt(int64(50), 10)
		}
		url = config.Conf.Clients["CoreData"].Url() + clients.ApiReadingRoute + "/device/" + device + "/" + limitUrl
	} else {
		url = config.Conf.Clients["CoreData"].Url() + clients.ApiReadingRoute
	}
	var readings []models.Reading

	err = request.Get(url, &readings)
	if err != nil {
		return
	}

	formatter := formatters.NewFormatter(readingTemplate, template.FuncMap{"DisplayDuration": utils.DisplayDuration})
	err = formatter.Write(readings)
	return
}
