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
	"os"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

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
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)

			for clientName, client := range config.Conf.Clients {
				resp, err := http.Get(client.Url() + clients.ApiPingRoute)
				if err != nil {
					fmt.Fprintf(w, "%s \t not connected\n", clientName)
				} else {
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Fprintf(w, "%s \t Unexpected error: %v\n", clientName, err)
					}
					if string(data) == "pong" {
						fmt.Fprintf(w, "%s \t OK\n", clientName)
					}
					resp.Body.Close()
				}
			}
			w.Flush()

		},
	}
	return cmd
}
