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
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

var age string
var slug string

func removeNotificationHandler(cmd *cobra.Command, args []string) (err error){
	if age == "" && slug == "" {
		return errors.New("at least one flag should be provided")
	} else if age != "" && slug != ""{
		return errors.New("age or slug flag should be provided, not both")
	}
	url, deletedBy := constructUrl()
	return request.DeletePrt(url, deletedBy)
}

// NewCommand returns the rm command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [--slug | --age]",
		Short: "Removes notification by slug or age",
		Long:  "Removes notifications by slug or age (the current notifications timestamp minus their last modification timestamp should be less than the age parameter)",
		RunE:   removeNotificationHandler,
	}
	cmd.Flags().StringVar(&age, "age", "", "Notification age (in milliseconds)")
	cmd.Flags().StringVarP(&slug, "slug", "s", "", "Meaningful, case insensitive identifier")
	return cmd
}

func constructUrl() (string, string) {
	url := config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute
	if age != "" {
		return fmt.Sprintf("%s/age/%s", url, age), age
	}
	return fmt.Sprintf("%s/slug/%s", url, slug), slug
}