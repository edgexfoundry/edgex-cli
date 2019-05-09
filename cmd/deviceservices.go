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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	models "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
)

type deviceServiceList struct {
	rd []models.DeviceProfile
}

// deviceservicesCmd represents the deviceservices command
var deviceservicesCmd = &cobra.Command{
	Use:   "deviceservices",
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
		for i, device := range deviceServiceList1.rd {
			fmt.Printf("%v\t%s\t%v\n", i, device.Name, device.Created)
		}
	},
}

func init() {
	rootCmd.AddCommand(deviceservicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceservicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceservicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
