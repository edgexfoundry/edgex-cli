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
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/responses"
)

// IssueReadCommand issues the specified read command referenced by the device name and command name.
// commandName: The name of the command to be executed;
// dsPushEvent: If set to yes, a successful read command will result in an event being pushed to the EdgeX system. Default value is no;
// dsReturnEvent: If set to no, there will be no event returned in the HTTP response. Default value is yes.
func (c Service) IssueReadCommand(deviceName string, commandName string, dsPushEvent string, dsReturnEvent string) (response *responses.EventResponse, err error) {
	client := c.getCommandClient()
	response, err = client.IssueGetCommandByName(context.Background(), deviceName, commandName, dsPushEvent, dsReturnEvent)
	return
}

// IssueWriteCommand issues the specified write command referenced by the device name, command name and settings.
// commandName: The name of the command to be executed;
// settings: It specifies the write command's request body, which provides the value(s) being written to the device.
func (c Service) IssueWriteCommand(deviceName string, commandName string, settings map[string]string) (response dtosCommon.BaseResponse, err error) {
	client := c.getCommandClient()
	response, err = client.IssueSetCommandByName(context.Background(), deviceName, commandName, settings)
	return
}

// ListAllCommands returns a sorted list of all available commands, optionally limiting the list by
// specifying the offset and limit parameters.
// offset: The number of items to skip. Default is 0;
// limit: The number of items to return (-1 will return all remaining items). Default is 50.
func (c Service) ListAllCommands(offset int, limit int) (response responses.MultiDeviceCoreCommandsResponse, err error) {
	client := c.getCommandClient()
	response, err = client.AllDeviceCoreCommands(context.Background(), offset, limit)
	return
}

// ListCommandsByDeviceName returns a sorted list of all available commands, filtered by a device name.
func (c Service) ListCommandsByDeviceName(deviceName string) (response responses.DeviceCoreCommandResponse, err error) {
	client := c.getCommandClient()
	response, err = client.DeviceCoreCommandsByDeviceName(context.Background(), deviceName)
	return
}

func (c Service) getCommandClient() interfaces.CommandClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewCommandClient(url)
}
