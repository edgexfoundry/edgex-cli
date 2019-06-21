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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	models "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
)

type deviceProfileList struct {
	rd []models.DeviceProfile
}

// NewCommand return the list profiles command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of device profiles",
		Long:  `Returns the list of device profiles currently in the core-metadata database.`,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := http.Get("http://localhost:48081/api/v1//deviceprofile")
			if err != nil {
				fmt.Println("An error occured. Is EdgeX running?")
				fmt.Println(err)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)

			deviceProfileList1 := deviceProfileList{}

			errjson := json.Unmarshal(data, &deviceProfileList1.rd)
			if errjson != nil {
				fmt.Println(errjson)
			}

			// TODO: Add commands and resources? They both are lists, so we need to think about how to display them
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "Profile ID", "Profile Name", "Created", "Modified", "Manufacturer", "Model")
			for _, device := range deviceProfileList1.rd {
				tCreated := time.Unix(device.Created/1000, 0)
				tModified := time.Unix(device.Modified/1000, 0)
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%v\t\n",
					device.Id,
					device.Name,
					humanDuration(time.Since(tCreated)),
					humanDuration(time.Since(tModified)),
					device.Manufacturer,
					device.Model,
				)
			}
			w.Flush()
		},
	}
	return cmd
}

// Taken from https://github.com/docker/go-units/blob/master/duration.go
func humanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds == 1 {
		return "1 second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours() + 0.5); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*2 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
}
