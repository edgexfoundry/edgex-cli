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
	"context"
	"fmt"
	"io/ioutil"
	netHttp "net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
)

// Service defines the hostname and port of a EdgeX microservice
type Service struct {

	// Host is the hostname
	Host string
	// Port number used by service
	Port int
}

//GetMetrics returns the metrics for this service.
func (c Service) GetMetrics() (result dtoCommon.Metrics, err error) {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	client := http.NewGeneralClient(url)
	response, err := client.FetchMetrics(context.Background())
	if err != nil {
		return result, err
	}

	return response.Metrics, nil

}

//GetVersionJSON returns the request URL and response for the 'version' endpoint.
func (c Service) GetVersionJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiVersionRoute)
}

//GetPingJSON returns the request URL and response for the 'ping' endpoint.
func (c Service) GetPingJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiPingRoute)
}

//GetConfigJSON returns the request URL and response for the 'config' endpoint.
func (c Service) GetConfigJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiConfigRoute)
}

//GetMetricsJSON returns the request URL and response for the 'metrics' endpoint.
func (c Service) GetMetricsJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiMetricsRoute)
}

//callEndpoint calls an endpoint on this service and returns the result and the URL used
func (c Service) callEndpoint(endpoint string) (string, string, error) {
	url := fmt.Sprintf("http://%s:%v%s", c.Host, c.Port, endpoint)

	resp, err := netHttp.Get(url)
	if err != nil {
		return "", "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", "", err
	}

	return string(data), url, nil

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
