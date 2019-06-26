package list

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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/edgexfoundry/edgex-cli/pkg/utils"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
)

// var rd []models.Device
type eventList struct {
	rd []models.Event
}

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all device services",
		Long:  `Return all device services sorted by id.`,
		Run: func(cmd *cobra.Command, args []string) {

			resp, err := http.Get("http://localhost:48080/api/v1/event")
			if err != nil {
				// handle error
				fmt.Println("An error occured. Is EdgeX running?")
				fmt.Println(err)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)

			eventList := eventList{}

			errjson := json.Unmarshal(data, &eventList.rd)
			if errjson != nil {
				fmt.Println(errjson)
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", "Event ID", "Device", "Origin", "Created", "Modified")
			for _, event := range eventList.rd {
				tCreated := time.Unix(event.Created/1000, 0)
				tModified := time.Unix(event.Modified/1000, 0)
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%s\t\n",
					event.ID,
					event.Device,
					event.Origin,
					utils.HumanDuration(time.Since(tCreated)),
					utils.HumanDuration(time.Since(tModified)),
				)
			}
			w.Flush()
		},
	}
	return cmd
}
