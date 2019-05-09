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

	"github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/spf13/cobra"
)

// var rd []models.Device
type deviceList struct {
	rd []models.Device
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "A list of all device services",
	Long: `Return all device services sorted by id. 
	Returns Internal Service Error (HTTP 500) for 
	unknown or unanticipated issues. Returns 
	LimitExceededException (HTTP 413) if the number 
	returned exceeds the max limit.`,
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := http.Get("http://localhost:48081/api/v1/device")
		if err != nil {
			// handle error
			fmt.Println("An error occured. Is EdgeX running?")
			fmt.Println(err)
		}
		defer resp.Body.Close()

		data, _ := ioutil.ReadAll(resp.Body)

		deviceList1 := deviceList{}

		errjson := json.Unmarshal(data, &deviceList1.rd)
		if errjson != nil {
			fmt.Println(errjson)
		}
		for i, device := range deviceList1.rd {
			fmt.Printf("%v\t%s\t%v\n", i, device.Name, device.Created)
		}

	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
