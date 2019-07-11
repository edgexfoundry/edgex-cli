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

package rm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// NewCommand returns the rm command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [device name]",
		Short: "Removes device by name",
		Long:  `Removes a device given its name.`,
		Run: func(cmd *cobra.Command, args []string) {

			client := &http.Client{}

			// Create request
			req, err := http.NewRequest("DELETE", "http://localhost:48081/api/v1/device/name/"+args[0], nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Fetch Request
			resp, errReq := client.Do(req)
			if errReq != nil {
				fmt.Println(errReq)
				return
			}
			defer resp.Body.Close()

			respBody, errrespBody := ioutil.ReadAll(resp.Body)
			if errrespBody != nil {
				fmt.Println(errrespBody)
				return
			}

			// Display Results

			if string(respBody) == "true" {
				fmt.Printf("Removed: %s\n", args[0])
			} else {
				fmt.Printf("Remove Unsuccessful!\n")
			}
		},
	}
	return cmd
}
