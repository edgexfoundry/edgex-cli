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

package rm

import (
	"errors"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

var slug string

// NewCommand returns remove subscription command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "rm [id | --slug]",
		Example: "subscription rm <id> \n" +
			"subscription rm --slug <slug>",
		Short: "Removes subscription by --slug or id.",
		Long:  `Removes a subscription given its slug or id.`,
		RunE:   removeSubscriptionHandler,
	}
	cmd.Flags().StringVarP(&slug, "slug", "s", "", "Delete Subscription by slug")
	return cmd
}

func removeSubscriptionHandler(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 && slug == "" {
		return errors.New("no subscription id/slug provided")
	}
	url := config.Conf.Clients["Notification"].Url() + clients.ApiSubscriptionRoute
	url, deletedBy := constructUrl(url, args)
	return request.DeletePrt(url, deletedBy)
}

func constructUrl(url string, args []string) (string, string) {
	if slug != "" {
		url = url + "/slug/" + slug
		return url, slug
	}
	url = url + "/" + args[0]
	return url, args[0]
}
