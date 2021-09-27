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
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
)

// RemoveEvents removes events, filtering either using a device name or
// a maximum event age, specified in milliseconds
func (c Service) RemoveEvents(device string, age int) error {
	client := c.getEventClient()

	if device != "" && age != 0 {
		return errors.New("either specify device name or event age, but not both")
	} else if device != "" {
		client.DeleteByDeviceName(context.Background(), device)
	} else if age != 0 {
		client.DeleteByAge(context.Background(), age)
	} else {
		return errors.New("event ID, device name or event age must be specified")
	}

	return nil
}

// CreateEvent creates a new event using the specified device, profile and source name and value type, with
// generating one or more sample readings
func (c Service) CreateEvent(deviceName string, profileName string, sourceName string, valueType string, numberOfReadings int) (string, error) {

	if numberOfReadings < 1 {
		return "", errors.New("the number of readings must be at least 1")

	}

	client := c.getEventClient()
	event := dtos.NewEvent(profileName, deviceName, sourceName)
	valueType = strings.Title(strings.ToLower(valueType))

	event.Readings = make([]dtos.BaseReading, numberOfReadings)
	for i := 0; i < numberOfReadings; i++ {
		var reading dtos.BaseReading
		var err error
		r64 := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
		var value interface{}
		switch valueType {
		case common.ValueTypeBool:
			value = (r64&1 == 0)
		case common.ValueTypeString:
			value = "Reading " + strconv.Itoa(i)
		case common.ValueTypeUint8:
			value = uint8(r64)
		case common.ValueTypeUint16:
			value = uint16(r64)
		case common.ValueTypeUint32:
			value = uint32(r64)
		case common.ValueTypeUint64:
			value = r64
		case common.ValueTypeInt8:
			value = int8(r64)
		case common.ValueTypeInt16:
			value = int16(r64)
		case common.ValueTypeInt32:
			value = int32(r64)
		case common.ValueTypeInt64:
			value = int64(r64)
		case common.ValueTypeFloat32:
			value = float32(r64) / 100
		case common.ValueTypeFloat64:
			value = float64(r64) / 100
		default:
			return "", errors.New("type must be one of [bool | string | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 ]")
		}

		reading, err = dtos.NewSimpleReading(profileName, deviceName, sourceName, valueType, value)

		if err != nil {
			return "", err
		}
		event.Readings[i] = reading

	}

	response, err := client.Add(context.Background(), requests.NewAddEventRequest(event))
	return response.Id, err
}

// ListAllEvents returns a sorted list of all available events, optionally limiting the list by
// specifying the offset and limit parameters.
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
func (c Service) ListAllEvents(offset, limit int) (events []dtos.Event, err error) {

	client := c.getEventClient()
	response, err := client.AllEvents(context.Background(), offset, limit)

	if err != nil {
		return nil, err
	}
	return response.Events, nil
}

// ListAllReadings returns a sorted list of all available readings, optionally limiting the list by
// specifying the offset and limit parameters.
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
func (c Service) ListAllReadings(offset, limit int) (events []dtos.BaseReading, err error) {

	client := c.getReadingClient()
	response, err := client.AllReadings(context.Background(), offset, limit)

	if err != nil {
		return nil, err
	}
	return response.Readings, nil
}

// CountEvents returns the number of events available, optionally filtered by a device name
func (c Service) CountEvents(device string) (response dtoCommon.CountResponse, err error) {

	client := c.getEventClient()

	if device != "" {
		response, err = client.EventCountByDeviceName(context.Background(), device)
	} else {
		response, err = client.EventCount(context.Background())
	}
	return
}

func (c Service) getEventClient() interfaces.EventClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewEventClient(url)

}

func (c Service) getReadingClient() interfaces.ReadingClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewReadingClient(url)

}
