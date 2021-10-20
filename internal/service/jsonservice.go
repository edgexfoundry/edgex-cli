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
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
)

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

//CountEventsByDeviceJSON returns the request URL and response for the 'event/count/device/name' endpoint.
func (c Service) CountEventsByDeviceJSON(device string) (json string, url string, err error) {
	endpoint := strings.Replace(common.ApiEventCountByDeviceNameRoute, "{name}", device, 1)
	return c.callEndpoint(endpoint)
}

//CountEventsJSON returns the request URL and response for the 'event/count' endpoint.
func (c Service) CountEventsJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiEventCountRoute)
}

//CountReadingsByDeviceJSON returns the request URL and response for the 'reading/count/device/name' endpoint.
func (c Service) CountReadingsByDeviceJSON(device string) (json string, url string, err error) {

	endpoint := strings.Replace(common.ApiReadingCountByDeviceNameRoute, "{name}", device, 1)

	return c.callEndpoint(endpoint)
}

//CountReadingsJSON returns the request URL and response for the 'reading/count' endpoint.
func (c Service) CountReadingsJSON() (json string, url string, err error) {
	return c.callEndpoint(common.ApiReadingCountRoute)
}

//ListReadingsJSON returns all readings
func (c Service) ListAllReadingsJSON(offset, limit int) (json string, urlString string, err error) {
	return c.getList(offset, limit, "", common.ApiAllReadingRoute)
}

//ListEventsJSON returns all events
func (c Service) ListAllEventsJSON(offset, limit int) (json string, urlString string, err error) {
	return c.getList(offset, limit, "", common.ApiAllEventRoute)
}

func (c Service) getList(offset, limit int, labels string, endpoint string) (json string, urlString string, err error) {
	if limit == -1 && offset == 0 && labels == "" {
		json, urlString, err = c.callEndpoint(endpoint)

	} else {
		var u *url.URL
		u, err = url.Parse(endpoint)
		if err != nil {
			return "", "", err
		}

		requestParams := url.Values{}
		requestParams.Set(common.Offset, strconv.Itoa(offset))
		requestParams.Set(common.Limit, strconv.Itoa(limit))
		if len(labels) > 0 {
			requestParams.Set(common.Labels, labels)
		}
		u.RawQuery = requestParams.Encode()
		json, urlString, err = c.callEndpoint(u.String())
	}
	return
}

//callEndpoint calls an endpoint on this service and returns the result and the URL used
func (c Service) callEndpoint(endpoint string) (string, string, error) {
	url := fmt.Sprintf("http://%s:%v%s", c.Host, c.Port, endpoint)

	resp, err := http.Get(url)
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
