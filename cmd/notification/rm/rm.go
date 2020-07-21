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
	"strconv"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/utils"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

const TimeUnitUsage = "Specify the time unit used with the --age flag\nList of possible values:\n" + utils.TimeUnitDescriptions

var age string
var unit string
var slug string

// NewCommand returns the rm command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [--slug | --age]",
		Short: "Removes notification by slug or age",
		Long:  "Removes notifications by slug or age (the current notifications timestamp minus their last modification timestamp should be greater than the age parameter)",
		RunE:  removeNotificationHandler,
	}
	cmd.Flags().StringVarP(&age, "age", "a", "", "Notification age (by default in milliseconds). To customize the unit use --unit flag")
	cmd.Flags().StringVarP(&slug, "slug", "s", "", "Meaningful, case insensitive identifier")
	cmd.Flags().StringVar(&unit, "unit", "ms", TimeUnitUsage)
	return cmd
}

func removeNotificationHandler(cmd *cobra.Command, args []string) (err error) {
	if age == "" && slug == "" {
		return errors.New("at least one flag should be provided")
	} else if age != "" && slug != "" {
		return errors.New("age or slug flag should be provided, not both")
	}

	url, deletedBy, err := constructUrl()
	if err != nil {
		return err
	}
	return request.DeletePrt(url, deletedBy)
}

func constructUrl() (string, string, error) {
	baseUrl := config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute
	url := fmt.Sprintf("%s/slug/%s", baseUrl, slug)
	if age != "" {
		if _, present := utils.TimeUnitsMap[unit]; !present {
			return "", "", errors.New("List of possible values:\n" + utils.TimeUnitDescriptions)
		}
		ageInt, err := strconv.ParseInt(age, 10, 64)
		if err != nil {
			return "", "", err
		}
		ageMilliseconds := utils.ConvertAgeToMillisecond(unit, ageInt)
		return fmt.Sprintf("%s/age/%v", baseUrl, ageMilliseconds), age + unit, nil
	}
	return url, slug, nil
}
