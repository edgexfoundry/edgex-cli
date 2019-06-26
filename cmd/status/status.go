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

	"github.com/spf13/cobra"
)

// NewCommand returns the status command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "[BETA]Checks the current status of each microservice.",
		Long: `Status
		
This command pings each edgex microservice and prints their status.
This command is not stable yet.
`,
		Run: func(cmd *cobra.Command, args []string) {

			// Maps of microservices pointing to their port numbers
			microservices := make(map[string]string)

			microservices["Core Data"] = "48080"
			microservices["Core Metadata"] = "48081"
			microservices["Core Command"] = "48082"
			microservices["Alerts & Notifications"] = "48060"
			microservices["Logging"] = "48061"
			microservices["Scheduling"] = "48085"
			microservices["Rules Engine"] = "48075"
			microservices["Client Registration"] = "48071"
			microservices["System Management"] = "48090"

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)

			for microservice, port := range microservices {
				resp, err := http.Get("http://localhost:" + port + "/api/v1/ping")
				if err != nil {
					fmt.Fprintf(w, "%s \t not connected\n", microservice)
				} else {
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Fprintf(w, "%s \t Unexpected error: %v\n", microservice, err)
					}
					if string(data) == "pong" {
						fmt.Fprintf(w, "%s \t OK\n", microservice)
					}
					resp.Body.Close()
				}
			}
			w.Flush()

		},
	}
	return cmd
}
