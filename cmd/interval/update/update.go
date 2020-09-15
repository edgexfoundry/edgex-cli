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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/scheduler"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

type IntervalFile struct {
	Intervals []models.Interval
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update interval",
		Long:  `Update the intervals described in the given TOML files.`,
		Run:   updateIntervalHandler,
	}
	return cmd
}

func updateIntervalHandler(cmd *cobra.Command, args []string) {
	for _, fname := range args {
		intervals, err := parseJson(fname)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			continue
		}

		for _, i := range intervals {
			updateInterval(i)
		}
	}
}

func updateInterval(n models.Interval) {
	url := config.Conf.Clients["Scheduler"].Url()
	client := scheduler.NewIntervalClient(
		local.New(url + clients.ApiIntervalRoute),
	)

	err := client.Update(context.Background(), n)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Interval successfully created: %s\n", n.Name)
	}
}

func parseJson(fname string) ([]models.Interval, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)

	fileContent := &IntervalFile{}
	err = json.Unmarshal([]byte(byteValue), &fileContent)
	if err != nil {
		return nil, err
	}
	return fileContent.Intervals, nil
}
