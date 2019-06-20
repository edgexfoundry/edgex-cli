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
			for i, profile := range deviceProfileList1.rd {
				fmt.Printf("%v\t%s\t%v\n", i, profile.Name, profile.Id)
			}
		},
	}
	return cmd
}
