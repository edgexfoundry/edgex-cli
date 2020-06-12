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

package status

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

const (
	NotConnected = "Not Connected"
	OK           = "Ok"
)

const statusTemplate = "Service Name\tURL\tStatus\n" +
	"{{range .}}" +
	"{{.Name}}\t{{.Url}}\t{{.Status}}\n" +
	"{{end}}"

// NewCommand returns the status command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Checks the current status of each microservice.",
		Long: `Status
		
This command pings each edgex microservice and prints their status.
This command is not stable yet.
`,
		Run: func(cmd *cobra.Command, args []string) {
			var clientStatuses []clientStatus
			for clientName, client := range config.Conf.Clients {
				resp, err := http.Get(client.Url() + clients.ApiPingRoute)
				if err != nil {
					clientStatuses = append(clientStatuses, clientStatus{clientName,client.Url(), NotConnected})
				} else {
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						clientStatuses = append(clientStatuses, clientStatus{clientName,client.Url(), fmt.Sprintf("%s \t Unexpected error: %v\n", clientName, err)})
					}
					if string(data) == "pong" {
						clientStatuses = append(clientStatuses, clientStatus{clientName,client.Url(), OK})
					}
					resp.Body.Close()
				}
			}
			formatters.NewHtmlTempleteFormatter("srvStatuses", statusTemplate, nil).Write(clientStatuses)
		},
	}
	return cmd
}

type clientStatus struct {
	Name string
	Url  string
	Status string
}
