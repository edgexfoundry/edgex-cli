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
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all subscriptions",
		Long:  `Return all Subscriptions`,
		RunE: func(cmd *cobra.Command, args []string) (err error){
			url:=config.Conf.Clients["Notification"].Url()+clients.ApiSubscriptionRoute
			var subscriptions []models.Subscription
			err = request.Get(url, &subscriptions)
			if err != nil {
				return
			}

			pw := viper.Get("writer").(io.Writer)
			w := new(tabwriter.Writer)
			w.Init(pw, 0, 8, 1, '\t', 0)
			_, err = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "Subscription ID", "Subscription Slug", "Origin")
			if err != nil {
				return
			}
			for _, subscription := range subscriptions {
				fmt.Fprintf(w, "%s\t%s\t%v\t\n",
					subscription.ID,
					subscription.Slug,
					subscription.Origin,
				)
			}

			err = w.Flush()
			if err != nil {
				return
			}
			return
		},
	}
	return cmd
}
