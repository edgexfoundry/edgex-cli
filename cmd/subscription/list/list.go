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
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/formatters"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

const subscriptionTemplete = "Subscription ID\tSlug\tReceiver\tDescription\tOrigin\n" +
	"{{range .}}" +
	"{{.ID}}\t{{.Slug}}\t{{.Receiver}}t{{.Description}}t{{.Origin}}\n" +
	"{{end}}"

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all subscriptions",
		Long:  `Return all Subscriptions`,
		RunE: listHandler,
	}
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	url := config.Conf.Clients["Notification"].Url() + clients.ApiSubscriptionRoute
	var subscriptions []models.Subscription
	err = request.Get(url, &subscriptions)
	if err != nil {
		return
	}

	formatter := formatters.NewFormatter("subscriptionList", subscriptionTemplete, nil)
	err = formatter.Write(subscriptions)
	return
}
