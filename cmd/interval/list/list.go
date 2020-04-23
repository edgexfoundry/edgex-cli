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
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var limit int32
var byID bool

func listHandler(cmd *cobra.Command, args []string) (err error){
	var url = config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute
	if len(args) > 0 {
		if byID {
			url += "/" + args[0]
		} else {
			url += "/name/" + args[0]
		}

	}
	var intervals []models.Interval
	err = client.ListHelper(url, &intervals)
	if err != nil {
		return
	}

	pw := viper.Get("writer").(io.WriteCloser)
	w := new(tabwriter.Writer)
	w.Init(pw, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", "Interval ID", "Name", "Start",
		"End", "Frequency", "Cron", "RunOnce")
	for _, interval := range intervals {
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
	w.Flush()
	return
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all intervals",
		Long:  `Return a list of all intervals.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:   listHandler,
	}
	cmd.Flags().BoolVar(&byID, "id", false, "By ID")
	return cmd
}
