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

type deviceServiceList struct {
	rd []models.DeviceService
}

// NewCommand returns the list device services command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists existing devices services",
		Long:  `Return the list fo current device services.`,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := http.Get("http://localhost:48081/api/v1/deviceservice")
			if err != nil {
				fmt.Println("An error occured. Is EdgeX running?")
				fmt.Println(err)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)

			deviceServiceList1 := deviceServiceList{}

			errjson := json.Unmarshal(data, &deviceServiceList1.rd)
			if errjson != nil {
				fmt.Println(errjson)
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "ID", "Name", "Last Connected", "Operating State")
			for _, device := range deviceServiceList1.rd {
				tCreated := time.Unix(0, device.LastConnected)
				elapsed := time.Since(tCreated)
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t\n",
					device.Id,
					device.Name,
					humanDuration(elapsed),
					device.OperatingState,
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
