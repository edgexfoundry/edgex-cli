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
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// See https://github.com/edgexfoundry/edgex-go/blob/master/api/openapi/v1/support-notifications.yaml
var limit int32
var start string
var end string
var slug string
var labels string
var sender string
var onlyNew bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all notifications",
		Long:  `Return a list of all notifications filtered by slug/sender/labels/start/end and limited by limit.`,
		Args:  cobra.MaximumNArgs(3),
		RunE:   listHandler,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
		TraverseChildren:           false,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	}

	cmd.Flags().Int32VarP(&limit, "limit", "l", 50, "Limit number of results")
	cmd.Flags().StringVar(&start,"start", "", "Filter results by start date")
	cmd.Flags().StringVar(&end, "end", "",  "Filter results by end date")
	cmd.Flags().StringVar(&slug,  "slug",  "",  "Filter results by slug")
	cmd.Flags().StringVar(&sender,  "sender",  "",  "Filter results by sender")
	cmd.Flags().StringVar(&labels,  "labels",  "",  "Filter results by labels")
	cmd.Flags().BoolVar(&onlyNew, "new", false, "Filter results by new")

	return cmd
}


func listHandler(cmd *cobra.Command, args []string) (err error){
	var url string
	// what determines the argument order?
	//limit could be skipped, both end and start could be specified, so max args == 3

	var startStr = cmd.Flag("start").Value.String()
	fmt.Println(startStr)
	url = config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute + "/"

	if onlyNew {
		url += "new"
	} else if slug != "" {
		url += "slug/" + slug
		limit = -1 // no limit with slug
	} else if labels != "" {
		url += "labels/" + labels
	} else if sender != "" {
		url += "sender/" + sender
	} else {
		if start != "" {
			url += "start/" + start
		}
		if end != "" {
			if strings.HasSuffix(url, "/") {
				url += "end/" + end
			} else { // start also specified?
				url += "/end/" + end
			}
		}
	}

	if limit > 0 {
		url = url + "/" + strconv.FormatInt(int64(limit), 10)
	}
	fmt.Printf ("*** URL ==  %s *** \n", url)
	var notifications []models.Notification
	err = request.Get(url, &notifications)
	if err != nil {
		return
	}

	pw := viper.Get("writer").(io.WriteCloser)
	w := new(tabwriter.Writer)
	w.Init(pw, 0, 8, 1, '\t', 0)
	//fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", "Interval ID", "Name", "Start",
	//	"End", "Frequency", "Cron", "RunOnce")
	for _, notification := range notifications {
		/*fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t\n",
			interval.ID,
			interval.Name,
			interval.Start,
			interval.End,
			interval.Frequency,
			interval.Cron,
			interval.RunOnce,
		)*/
		fmt.Fprintf(w, "%s\n", notification)
	}
	w.Flush()
	return
}


