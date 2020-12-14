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

package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/edgexfoundry/edgex-cli/cmd/interval/add"
	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/scheduler"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var name string
var file string

// NewCommand returns the update Interval command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update interval",
		Long: "Update interval(s) described in the given JSON file or use the interactive mode enabled by providing " +
			"name of existing interval. \n" +
			"Parameters description: \n" +
			fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", add.IntervalStartUsage, add.IntervalEndUsage, add.RunOnceIntervalUsage, add.CronIntervalUsage, add.FrequencyIntervalUsage),
		RunE: intervalHandler,
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Interval name. Interval with given name is loaded into default editor, ready to be customized")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Json file containing interval configuration to update")
	return cmd
}

func intervalHandler(cmd *cobra.Command, args []string) error {
	if name != "" && file != "" {
		return errors.New("Interval could be updated by providing a file, or by specifying interval name to be updated using interactive mode. ")
	}

	if name == "" && file == "" {
		return errors.New("Please, provide file or interval name ")
	}

	if file != "" {
		return updateIntervalFromFile(cmd)
	}

	updatedInterval, err := parseInterval(cmd, name)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute)
	err = scheduler.NewIntervalClient(client).Update(cmd.Context(), updatedInterval)
	if err != nil {
		return err
	}

	return nil
}

//parseInterval loads a Interval to be updated and open a default editor for customization
func parseInterval(cmd *cobra.Command, name string) (models.Interval, error) {
	var err error
	client := local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute)
	i, err := scheduler.NewIntervalClient(client).IntervalForName(cmd.Context(), name)
	if err != nil {
		return models.Interval{}, err
	}

	updatedIntervalBytes, err := editor.OpenInteractiveEditor(i, add.BaseIntervalTemp, template.FuncMap{
		"lastElem": editor.IsLastElementOfSlice,
	})

	if err != nil {
		return models.Interval{}, err
	}
	var updatedInterval models.Interval
	err = json.Unmarshal(updatedIntervalBytes, &updatedInterval)
	if err != nil {
		return models.Interval{}, errors.New("Unable to execute the command. The provided information is invalid: " + err.Error())
	}
	return updatedInterval, err
}

func updateIntervalFromFile(cmd *cobra.Command) error {
	intervals, err := LoadDSFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute)
	for _, ds := range intervals {
		err = scheduler.NewIntervalClient(client).Update(cmd.Context(), ds)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

//loadJsonFile could read a file that contains single Interval or list of Intervals
func LoadDSFromFile(filePath string) ([]models.Interval, error) {
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
	var ds models.Interval
	err = json.Unmarshal(file, &ds)
	if err != nil {
		//check if the file contains list of Interval
		err = json.Unmarshal(file, &intervals)
		if err != nil {
			return nil, err
		}
	} else {
		intervals = append(intervals, ds)
	}
	return intervals, nil
}
