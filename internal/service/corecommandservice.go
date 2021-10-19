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
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/responses"
	"strings"
)

func (c Service) IssueReadCommand(deviceName string, commandName string, dsPushEvent string, dsReturnEvent string) (response *responses.EventResponse, err error) {
	client := c.getCommandClient()
	response, err = client.IssueGetCommandByName(context.Background(), deviceName, commandName, dsPushEvent, dsReturnEvent)
	return
}

func (c Service) IssueWriteCommand(deviceName string, commandName string, settings map[string]string) (response dtosCommon.BaseResponse, err error) {
	client := c.getCommandClient()
	response, err = client.IssueSetCommandByName(context.Background(), deviceName, commandName, settings)
	return
}

func (c Service) ListAllCommands(offset int, limit int) (response responses.MultiDeviceCoreCommandsResponse, err error) {
	client := c.getCommandClient()
	response, err = client.AllDeviceCoreCommands(context.Background(), offset, limit)
	return
}

func (c Service) ListCommandsByDeviceName(deviceName string) (response responses.DeviceCoreCommandResponse, err error) {
	client := c.getCommandClient()
	response, err = client.DeviceCoreCommandsByDeviceName(context.Background(), deviceName)
	return
}

func (c Service) GetReadEndpoint(deviceName string, commandName string, dsPushEvent string, dsReturnEvent string) string {
	url := c.getEndpointUrl(common.ApiDeviceNameCommandNameRoute)
	replacer := strings.NewReplacer("{name}", deviceName, "{command}", commandName)
	return replacer.Replace(url) + "?ds-pushevent=" + dsPushEvent + "&ds-returnevent=" + dsReturnEvent
}

func (c Service) GetWriteEndpoint(deviceName string, commandName string, requestBody string) string {
	url := c.getEndpointUrl(common.ApiDeviceNameCommandNameRoute)
	replacer := strings.NewReplacer("{name}", deviceName, "{command}", commandName)
	return replacer.Replace(url) + " -d '" + requestBody + "'"
}

func (c Service) GetListAllEndpoint() string {
	return c.getEndpointUrl(common.ApiAllDeviceRoute)
}

func (c Service) GetListByDeviceEndpoint(deviceName string) string {
	url := c.getEndpointUrl(common.ApiDeviceByNameRoute)
	replacer := strings.NewReplacer("{name}", deviceName)
	return replacer.Replace(url)
}

func (c Service) getEndpointUrl(endpoint string) string {
	return fmt.Sprintf("http://%s:%v%s", c.Host, c.Port, endpoint)
}

func (c Service) getCommandClient() interfaces.CommandClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewCommandClient(url)
}
