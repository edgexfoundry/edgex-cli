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
	"io"
	"io/ioutil"
	"net/http"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type intervalList struct {
	intervals []models.Interval
}

var limit int32
var byID bool

func listHandler(cmd *cobra.Command, args []string) {
	var url string = "http://" + viper.GetString("Host") + ":48085/api/v1/"
	if len(args) > 0 {
		if byID {
			url += "interval/" + args[0]
		} else {
			url += "interval/name/" + args[0]
		}

	} else {
		url += "interval"
	}
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Println("An error occurred. Is EdgeX running?")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	intervalList := intervalList{}
	var errjson error
	var interval models.Interval
	if len(args) > 0 {
		errjson = json.Unmarshal(data, &interval)
	} else {
		errjson = json.Unmarshal(data, &intervalList.intervals)
	}
	if errjson != nil {
		if string(data) == "Error, exceeded the max limit as defined in config" {
			fmt.Println("The number of intervals to be returned exceeds the MaxResultCount limit defined in configuration.toml")
		}
		fmt.Println(errjson)
		return
	}

	pw := viper.Get("writer").(io.WriteCloser)
	w := new(tabwriter.Writer)
	w.Init(pw, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", "Interval ID", "Name", "Start",
		"End", "Frequency", "Cron", "RunOnce")
	if len(args) > 0 {

		printIntervalDetails(w, &interval)
	} else {
		for _, interval := range intervalList.intervals {
			printIntervalDetails(w, &interval)
		}
	}

	w.Flush()
}

func printIntervalDetails(w *tabwriter.Writer, interval *models.Interval) {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t\n",
		interval.ID,
		interval.Name,
		interval.Start,
		interval.End,
		interval.Frequency,
		interval.Cron,
		interval.RunOnce,
	)
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all intervals",
		Long:  `Return a list of all intervals.`,
		Args:  cobra.MaximumNArgs(1),
		Run:   listHandler,
	}
	cmd.Flags().BoolVar(&byID, "id", false, "By ID")
	return cmd
}
