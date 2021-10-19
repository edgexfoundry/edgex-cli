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
	"github.com/edgexfoundry/edgex-cli/internal/config"
	"github.com/edgexfoundry/edgex-cli/internal/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/spf13/cobra"
)

var verbose, metadata, data, command, notifications, scheduler, json bool

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

func addFormatFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&json, "json", "j", false, "show the raw JSON response")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show verbose/debug output")

}

func addStandardFlags(cmd *cobra.Command) {
	addFormatFlags(cmd)
	cmd.Flags().BoolVarP(&data, "data", "d", false, "use core-data service endpoint")
	cmd.Flags().BoolVarP(&command, "command", "c", false, "use core-command service endpoint")
	cmd.Flags().BoolVarP(&metadata, "metadata", "m", false, "use core-metadata service endpoint")
	cmd.Flags().BoolVarP(&scheduler, "scheduler", "s", false, "use support-scheduler service endpoint")
	cmd.Flags().BoolVarP(&notifications, "notifications", "n", false, "use support-notifications service endpoint")

}
