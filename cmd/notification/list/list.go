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
	"html/template"
	"strconv"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

// See https://github.com/edgexfoundry/edgex-go/blob/master/api/openapi/v1/support-notifications.yaml
var limit int32
var start string
var end string
var slug string
var labels string
var sender string
var onlyNew bool

const notificationTemplate = "Notification ID\tSlug\tSender\tStatus\tSeverity\tCategory\tContent\tLabels\tCreated\tModified\n" +
	"{{range .}}" +
	"{{.ID}}\t{{.Slug}}\t{{.Sender}}\t{{.Status}}\t{{.Severity}}\t{{.Category}}\t{{.Content}}\t{{.Labels}}\t{{DisplayDuration .Created}}\t{{DisplayDuration .Modified}}\n" +
	"{{end}}"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "list",
		Short:                      "A list of all notifications",
		Long:                       `Return a list of all notifications filtered by slug/sender/labels/start/end/new and limited by limit. Defaults to new notifications.`,
		Args:                       cobra.MaximumNArgs(3),
		RunE:                       listHandler,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	}

	cmd.Flags().Int32VarP(&limit, "limit", "l", 50, "Limit number of results")
	cmd.Flags().StringVar(&start, "start", "", "Filter results by start date")
	cmd.Flags().StringVar(&end, "end", "", "Filter results by end date")
	cmd.Flags().StringVar(&slug, "slug", "", "Filter results by slug")
	cmd.Flags().StringVar(&sender, "sender", "", "Filter results by sender")
	cmd.Flags().StringVar(&labels, "labels", "", "Filter results by labels")
	cmd.Flags().BoolVar(&onlyNew, "new", false, "Filter results by new")

	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	var url string
	multi := true
	url = config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute

	// For slug and id based retrieval, response will be a single item at most
	if slug != "" {
		url += "/slug/" + slug
		multi = false // no limit with slug
	} else if len(args) == 1 {
		// notification id provided
		multi = false
		url = url + "/" + args[0]
	} else if onlyNew {
		url += "/new"
	} else if labels != "" {
		url += "/labels/" + labels
	} else if sender != "" {
		url += "/sender/" + sender
	} else if start != "" {
		url += "/start/" + start
		// end could also be specified
		if end != "" {
			url += "/end/" + end
		}
	} else if end != "" {
		url += "/end/" + end
	} else { // default behavior whgen no flags specified
		url += "/new"
	}

	if multi {
		url = url + "/" + strconv.FormatInt(int64(limit), 10)
	}
	var notifications []models.Notification
	var aNotification models.Notification
	if !multi {
		err = request.Get(url, &aNotification)
	} else {
		err = request.Get(url, &notifications)
	}
	if err != nil {
		return
	}
	if !multi { // to use the same display code
		notifications = []models.Notification{aNotification}
	}
	formatter := formatters.NewFormatter(notificationTemplate, template.FuncMap{"DisplayDuration": utils.DisplayDuration})
	err = formatter.Write(notifications)
	return
}
