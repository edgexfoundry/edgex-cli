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
	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

const unitValues = "List of possible values:\n" +
	"ms - milliseconds\n" +
	"s - seconds\n" +
	"h - hours\n" +
	"d - days\n"

const unitUsage = "Specify the time unit used with the --age flag\n" + unitValues

var age string
var unit string
var slug string

//empty struct has width of zero.  It occupies zero bytes of storage
var units = map[string]struct{}{"ms": {}, "s": {}, "h": {}, "d": {}, "m": {}}

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
	cmd.Flags().StringVar(&unit, "unit", "ms", unitUsage)
	return cmd
}

func removeNotificationHandler(cmd *cobra.Command, args []string) (err error) {
	if age == "" && slug == "" {
		return errors.New("at least one flag should be provided")
	} else if age != "" && slug != "" {
		return errors.New("age or slug flag should be provided, not both")
	}
	if _, present := units[unit]; !present {
		return errors.New(unitValues)
	}
	ageInt, err := strconv.ParseInt(age, 10, 64)
	if err != nil {
		return errors.New("age is not an number")
	}
	url, deletedBy := constructUrl(ageInt)
	return request.DeletePrt(url, deletedBy)
}

func convertAgeToMillisecond(age int64) int64 {
	var ageMilliseconds int64
	switch unit {
	case "ms":
		ageMilliseconds = age
	case "s":
		ageMilliseconds = age * 1000
	case "m":
		ageMilliseconds = age * 60 * 1000
	case "h":
		ageMilliseconds = age * 60 * 60 * 1000
	case "d":
		ageMilliseconds = age * 24 * 60 * 60 * 1000
	}
	return ageMilliseconds
}

func constructUrl(ageInt int64) (string, string) {
	url := config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute
	if age != "" {
		ageMilliseconds := convertAgeToMillisecond(ageInt)
		return fmt.Sprintf("%s/age/%v", url, ageMilliseconds), age + unit
	}
	return fmt.Sprintf("%s/slug/%s", url, slug), slug
}
