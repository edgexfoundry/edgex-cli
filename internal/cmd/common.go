/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/edgexfoundry/edgex-cli/internal/config"
	"github.com/edgexfoundry/edgex-cli/internal/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/spf13/cobra"
)

var verbose, metadata, data, command, notifications, scheduler, json bool
var limit, offset int
var labels string

func getSelectedServiceKey() string {
	if metadata {
		return common.CoreMetaDataServiceKey
	} else if data {
		return common.CoreDataServiceKey
	} else if command {
		return common.CoreCommandServiceKey
	} else if notifications {
		return common.SupportNotificationsServiceKey
	} else if scheduler {
		return common.SupportSchedulerServiceKey
	} else {
		return ""
	}
}

func getSupportSchedulerService() service.Service {
	return config.GetCoreService(common.SupportSchedulerServiceKey)
}

func getSupportNotificationsService() service.Service {
	return config.GetCoreService(common.SupportNotificationsServiceKey)
}

func getCoreMetaDataService() service.Service {
	return config.GetCoreService(common.CoreMetaDataServiceKey)
}

func getCoreDataService() service.Service {
	return config.GetCoreService(common.CoreDataServiceKey)
}

func getCoreCommandService() service.Service {
	return config.GetCoreService(common.CoreCommandServiceKey)
}

func getSelectedServices() map[string]service.Service {
	key := getSelectedServiceKey()
	if key == "" {
		if json {
			key = common.CoreMetaDataServiceKey
		} else {
			return config.GetCoreServices()
		}
	}
	return map[string]service.Service{key: config.GetCoreService(key)}

}

func addVerboseFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
}

func addFormatFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&json, "json", "j", false, "Show the raw JSON response")

}

func addLimitOffsetFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "The number of items to return. Specifying -1 will return all remaining items")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "The number of items to skip")
}

func addLabelsFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&labels, "labels", "", "", "Comma-delimited list of user-defined labels")
}

func getLabels() []string {
	var aLabels []string
	if len(labels) > 0 {
		aLabels = strings.Split(labels, ",")
	}
	return aLabels
}

func validateAdminState(adminState string) error {
	if !(adminState == models.Locked || adminState == models.Unlocked) {
		return fmt.Errorf("admin state should be %s or %s", models.Locked, models.Unlocked)
	}
	return nil
}

func validateOperatingState(operState string) error {
	if !(operState == models.Up || operState == models.Down || operState == models.Unknown) {
		return fmt.Errorf("operating state should be one of %s,%s or %s", models.Up, models.Down, models.Unknown)
	}
	return nil
}

func addStandardFlags(cmd *cobra.Command) {
	addFormatFlags(cmd)
	cmd.Flags().BoolVarP(&data, "data", "d", false, "use core-data service endpoint")
	cmd.Flags().BoolVarP(&command, "command", "c", false, "use core-command service endpoint")
	cmd.Flags().BoolVarP(&metadata, "metadata", "m", false, "use core-metadata service endpoint")
	cmd.Flags().BoolVarP(&scheduler, "scheduler", "s", false, "use support-scheduler service endpoint")
	cmd.Flags().BoolVarP(&notifications, "notifications", "n", false, "use support-notifications service endpoint")

}

func getRFC822Time(t int64) string {
	if t == 0 {
		return "0"
	} else {
		return time.Unix(0, t*int64(time.Millisecond)).Format(time.RFC822)
	}
}

func getMillisTimestampFromRFC822Time(t string) (int64, error) {
	result, err := time.Parse(time.RFC822, t)
	if err != nil {
		return 0, err
	}
	return result.UnixNano() / int64(time.Millisecond), nil
}
