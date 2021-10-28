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

package service

import (
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
)

// Service defines the hostname and port of a EdgeX microservice
type Service struct {

	// Host is the hostname
	Host string
	// Port number used by service
	Port int
}

func (c Service) GetCommonClient() interfaces.CommonClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewCommonClient(url)
}

func (c Service) GetCommandClient() interfaces.CommandClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewCommandClient(url)
}

func (c Service) GetEventClient() interfaces.EventClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewEventClient(url)
}

func (c Service) GetReadingClient() interfaces.ReadingClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewReadingClient(url)
}

func (c Service) GetProvisionWatcherClient() interfaces.ProvisionWatcherClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewProvisionWatcherClient(url)
}

func (c Service) GetDeviceClient() interfaces.DeviceClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceClient(url)
}

func (c Service) GetDeviceServiceClient() interfaces.DeviceServiceClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceServiceClient(url)
}

func (c Service) GetDeviceProfileClient() interfaces.DeviceProfileClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceProfileClient(url)
}

func (c Service) GetNotificationClient() interfaces.NotificationClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewNotificationClient(url)
}
