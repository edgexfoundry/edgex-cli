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

package config

import (
	"github.com/edgexfoundry/edgex-cli/internal/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
)

var configuration HostConfiguration

type HostConfiguration struct {
	// CoreServices is a map of the core EdgeX microservices
	CoreServices Services
}

type Services map[string]service.Service

// GetCoreService returns the configuration of a core service
func GetCoreService(name string) service.Service {
	return configuration.CoreServices[name]
}

// GetCoreServices returns a map of the core EdgeX microservices
func GetCoreServices() Services {
	return configuration.CoreServices
}

func init() {
	configuration.CoreServices = Services{
		common.CoreMetaDataServiceKey: {
			Host: "localhost",
			Port: 59881,
		},
		common.CoreDataServiceKey: {
			Host: "localhost",
			Port: 59880,
		},
		common.CoreCommandServiceKey: {
			Host: "localhost",
			Port: 59882,
		},
		common.SupportSchedulerServiceKey: {
			Host: "localhost",
			Port: 59861,
		},
		common.SupportNotificationsServiceKey: {
			Host: "localhost",
			Port: 59860,
		},
	}

}
