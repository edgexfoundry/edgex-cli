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

package rm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// NewCommand returns the rm device service command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [device service name]",
		Short: "Removes device service by name",
		Long:  `Removes a device service from the core-metadata DB.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args[0])

			client := &http.Client{}

			// Create request
			req, err := http.NewRequest("DELETE", "http://localhost:48081/api/v1/deviceservice/name/"+args[0], nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Fetch Request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Display Results
			fmt.Println("response Status : ", resp.Status)
			fmt.Println("response Headers : ", resp.Header)
			fmt.Println("response Body : ", string(respBody))
		},
	}
	return cmd
}
