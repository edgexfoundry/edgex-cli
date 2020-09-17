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

package add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/scheduler"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"io/ioutil"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var interactiveMode bool
var name string
var start string
var end string
var cron string
var runOnce bool
var frequency string

var file string

const IntervalTempl = `[{{range $interval := .}}` + baseIntervalTemp +
	`{{end}}]`

const baseIntervalTemp = `{
	"Name" : "{{.Name}}",
	"start" : "{{.Start}}",
	"end" : "{{.End}}",
	"cron" : "{{.Cron}}",
	"runOnce" : {{.RunOnce}},
	"frequency" : "{{.Frequency}}"
}
`

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add interval",
		Long:  `Create interval(s) described in the given JSON file or use the interactive mode with additional flags.`,
		RunE:  newIntervalHandler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default " +
		"editor to customize the Interval information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Interval name")
	cmd.Flags().StringVarP(&start, "start", "s", "", "Interval Start time in format YYYYMMDD'T'HHmmss")
	cmd.Flags().StringVarP(&end, "end", "e", "", "Interval End time in format YYYYMMDD'T'HHmmss")
	cmd.Flags().StringVarP(&cron, "cron", "c", "", "Styled regular expression indicating " +
		"how often the action under interval should occur. Use either runOnce, frequency or cron and not all.")
	cmd.Flags().BoolVar(&runOnce, "runOnce",false, "runOnce - boolean indicating that this interval " +
		"runs one time - at the time indicated by the start")
	cmd.Flags().StringVar(&frequency, "frequency", "", "Interval frequency - how frequently should the\n" +
		"event occur. It is a sequence of decimal numbers, each with optional fraction and a unit suffix," +
		"such as \"300ms\", \"1.5h\" or \"2h45m\" or \"1h15m30s10us9ns\".\n" +
		"// Valid time units are:" +
		" \"ns\" (nanoseconds)," +
		"\"us\" (or \"µs\" for microseconds)," +
		"\"ms\"(for milliseconds), " +
		"\"s\"(seconds)," +
		"\"m\"(minutes)," +
		" \"h\"(hours)")

	cmd.Flags().StringVarP(&file, "file", "f", "", "Json file containing interval(s) configuration")
	return cmd
}

func newIntervalHandler(cmd *cobra.Command, args []string) error {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	if file != "" {
		return createIntervalsFromFile()
	}

	intervals, err := parseInterval(interactiveMode)
	if err != nil {
		return err
	}

	client := scheduler.NewIntervalClient(
		local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute),
	)
	for _, interval := range intervals {
		_, err = client.Add(context.Background(), &interval)
		if err != nil {
			return err
		}
	}
	return nil
}

func populateInterval(intervals *[]models.Interval) {
	i := models.Interval{}
	i.Name = name
	i.Start = start
	i.End = end
	i.Cron = cron
	i.RunOnce = runOnce
	i.Frequency = frequency
	*intervals = append(*intervals, i)
}

func createIntervalsFromFile() error {
	intervals, err := loadIntervalFromFile(file)
	if err != nil {
		return err
	}

	client := scheduler.NewIntervalClient(
		local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute),
	)
	for _, i := range intervals {
		_, err = client.Add(context.Background(), &i)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

func parseInterval(interactiveMode bool) ([]models.Interval, error) {
	//parse Intervals based on interactive mode and the other provided flags
	var err error
	var intervals []models.Interval
	populateInterval(&intervals)

	var updatedIntervalBytes []byte
	if interactiveMode {
		updatedIntervalBytes, err = editor.OpenInteractiveEditor(intervals, IntervalTempl, nil)
	} else {
		updatedIntervalBytes, err = json.Marshal(intervals)
	}
	if err != nil {
		return nil, err
	}

	var updatedIntervals []models.Interval
	err = json.Unmarshal(updatedIntervalBytes, &updatedIntervals)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is invalid: " + err.Error())
	}
	return updatedIntervals, err
}

//loadIntervalFromFile could read a file that contains single Interval or list of Intervals
func loadIntervalFromFile(filePath string) ([]models.Interval, error){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var intervals []models.Interval

	//check if the file contains just one Interval
	var i models.Interval
	err = json.Unmarshal(file, &i)
	if err != nil {
		//check if the file contains list of Interval
		err = json.Unmarshal(file, &intervals)
		if err != nil {
			return nil, err
		}
	} else {
		intervals = append(intervals, i)
	}
	return intervals, nil
}
