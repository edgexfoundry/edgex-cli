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
	gohttp "net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
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

//GetVersionJSON returns the request URL and response for the 'version' endpoint.
func (c Service) GetVersionJSON() (json string, url string, err error) {
	json, url, err = c.callEndpoint(common.ApiVersionRoute)
	return
}

//GetPingJSON returns the request URL and response for the 'ping' endpoint.
func (c Service) GetPingJSON() (json string, url string, err error) {
	json, url, err = c.callEndpoint(common.ApiPingRoute)
	return
}

//GetConfigJSON returns the request URL and response for the 'config' endpoint.
func (c Service) GetConfigJSON() (json string, url string, err error) {
	json, url, err = c.callEndpoint(common.ApiConfigRoute)
	return
}

//GetMetricsJSON returns the request URL and response for the 'metrics' endpoint.
func (c Service) GetMetricsJSON() (json string, url string, err error) {
	json, url, err = c.callEndpoint(common.ApiMetricsRoute)
	return
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

//callEndpoint calls an endpoint on this service and returns the result and the URL used
func (c Service) callEndpoint(endpoint string) (string, string, error) {
	url := fmt.Sprintf("http://%s:%v%s", c.Host, c.Port, endpoint)

	resp, err := gohttp.Get(url)
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
