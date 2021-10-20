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

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
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
