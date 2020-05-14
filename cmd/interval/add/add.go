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
	"fmt"
	"io/ioutil"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type IntervalFile struct {
	Intervals []models.Interval
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add interval",
		Long:  `Create the intervals described in the given TOML files.`,
		Run:   addIntervalHandler,
	}
	return cmd
}

func addIntervalHandler(cmd *cobra.Command, args []string) {
	for _, fname := range args {
		intervals, err := parseToml(fname)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			continue
		}
		for _, i := range intervals {
			request.Post(config.Conf.Clients["Scheduler"].Url()+clients.ApiIntervalRoute, &i)
		}
	}
}

func parseToml(fname string) ([]models.Interval, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid TOML")
		}
	}()
	var intervalFile = &IntervalFile{}
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(file, intervalFile)
	if err != nil {
		return nil, err
	}
	return intervalFile.Intervals, nil
}
