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

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

// ListAllDevices returns a list of all available devices
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
// labels: Array of label names
func (c Service) ListAllDevices(offset, limit int, labels []string) ([]dtos.Device, error) {

	client := c.getDeviceClient()
	response, err := client.AllDevices(context.Background(), labels, offset, limit)

	if err != nil {
		return nil, err
	}

	return response.Devices, nil
}

// ListAllDeviceServices returns a list of all available device services
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
// labels: Array of label names
func (c Service) ListAllDeviceServices(offset, limit int, labels []string) ([]dtos.DeviceService, error) {

	client := c.getDeviceServiceClient()
	response, err := client.AllDeviceServices(context.Background(), labels, offset, limit)

	if err != nil {
		return nil, err
	}
	return response.Services, nil
}

// ListAllDeviceProfiles returns a list of all available device profiles
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
// labels: Array of label names
func (c Service) ListAllDeviceProfiles(offset, limit int, labels []string) ([]dtos.DeviceProfile, error) {

	client := c.getDeviceProfileClient()
	response, err := client.AllDeviceProfiles(context.Background(), labels, offset, limit)

	if err != nil {
		return nil, err
	}
	return response.Profiles, nil
}

// ListAllProvisionWatchers returns a list of all available provision watchers
// offset: The number of items to skip. Default is 0.
// limit: The number of items to return (-1 will return all remaining items).
// labels: Array of label names
func (c Service) ListAllProvisionWatchers(offset, limit int, labels []string) ([]dtos.ProvisionWatcher, error) {

	client := c.getProvisionWatcherClient()
	response, err := client.AllProvisionWatchers(context.Background(), labels, offset, limit)

	if err != nil {
		return nil, err
	}
	return response.ProvisionWatchers, nil
}

// AddDeviceProfile adds a new device profile
func (c Service) AddDeviceProfile(name, description, manufacturer, model string, labels []string,
	resources []dtos.DeviceResource, commands []dtos.DeviceCommand) (*common.BaseResponse, error) {
	client := c.getDeviceProfileClient()

	var req = requests.NewDeviceProfileRequest(dtos.DeviceProfile{
		Name:            name,
		Description:     description,
		Manufacturer:    manufacturer,
		Model:           model,
		Labels:          labels,
		DeviceResources: resources,
		DeviceCommands:  commands,
	})

	response, err := client.Add(context.Background(), []requests.DeviceProfileRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response[0].BaseResponse, nil
}

// AddDeviceService adds a new device service
func (c Service) AddDeviceService(name, description string, baseAddress string, adminState string, labels []string) (*common.BaseResponse, error) {
	client := c.getDeviceServiceClient()

	err := validateAdminState(adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewAddDeviceServiceRequest(dtos.DeviceService{
		Name:        name,
		Description: description,
		Labels:      labels,
		BaseAddress: baseAddress,
		AdminState:  adminState,
	})

	response, err := client.Add(context.Background(), []requests.AddDeviceServiceRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response[0].BaseResponse, nil
}

// AddDevice provisions a new device
func (c Service) AddDevice(name, description, service, profile, adminState, operState string, labels []string, location string, protocols map[string]dtos.ProtocolProperties) (*common.BaseResponse, error) {
	client := c.getDeviceClient()

	err := validateOperatingState(operState)
	if err != nil {
		return nil, err
	}

	err = validateAdminState(adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewAddDeviceRequest(dtos.Device{
		Name:           name,
		Description:    description,
		ServiceName:    service,
		ProfileName:    profile,
		AdminState:     adminState,
		OperatingState: operState,
		Labels:         labels,
		Location:       location,
		AutoEvents:     nil,
		Protocols:      protocols,
	})
	response, err := client.Add(context.Background(), []requests.AddDeviceRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response[0].BaseResponse, nil
}

// AddProvisionWatcher adds a new provision watcher
func (c Service) AddProvisionWatcher(name, service, profile, adminState string, labels []string, identifiers map[string]string) (*common.BaseResponse, error) {
	client := c.getProvisionWatcherClient()

	err := validateAdminState(adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewAddProvisionWatcherRequest(dtos.ProvisionWatcher{
		Name:                name,
		ServiceName:         service,
		ProfileName:         profile,
		AdminState:          adminState,
		Labels:              labels,
		Identifiers:         identifiers,
		BlockingIdentifiers: nil,
	})
	response, err := client.Add(context.Background(), []requests.AddProvisionWatcherRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response[0].BaseResponse, nil
}

// GetDeviceProfileByName returns a named device
func (c Service) GetDeviceProfileByName(name string) (*dtos.DeviceProfile, error) {
	client := c.getDeviceProfileClient()

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	response, err := client.DeviceProfileByName(context.Background(), name)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response.Profile, nil
}

// GetDeviceByName returns a named device
func (c Service) GetDeviceByName(name string) (*dtos.Device, error) {
	client := c.getDeviceClient()

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	response, err := client.DeviceByName(context.Background(), name)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response.Device, nil
}

// GetDeviceService returns a named device service
func (c Service) GetDeviceServiceByName(deviceService string) (*dtos.DeviceService, error) {
	client := c.getDeviceServiceClient()

	err := validateName(deviceService)
	if err != nil {
		return nil, err
	}

	response, err := client.DeviceServiceByName(context.Background(), deviceService)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response.Service, nil
}

// GetDeviceService returns a named device service
func (c Service) GetProvisionWatcherByName(name string) (*dtos.ProvisionWatcher, error) {
	client := c.getProvisionWatcherClient()

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	response, err := client.ProvisionWatcherByName(context.Background(), name)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &response.ProvisionWatcher, nil
}

// RemoveDevice removes a named device
func (c Service) RemoveDevice(deviceService string) (*common.BaseResponse, error) {
	client := c.getDeviceClient()
	err := validateName(deviceService)

	if err != nil {
		return nil, err
	}

	response, err := client.DeleteDeviceByName(context.Background(), deviceService)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &response, nil
}

// RemoveDeviceService removes a named device service
func (c Service) RemoveDeviceService(deviceService string) (*common.BaseResponse, error) {
	client := c.getDeviceServiceClient()

	err := validateName(deviceService)
	if err != nil {
		return nil, err
	}

	response, err := client.DeleteByName(context.Background(), deviceService)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &response, nil
}

// RemoveDeviceProfile removes a named device profile
func (c Service) RemoveDeviceProfile(name string) (*common.BaseResponse, error) {
	client := c.getDeviceProfileClient()

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	response, err := client.DeleteByName(context.Background(), name)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &response, nil
}

// RemoveDeviceProfile removes a named device profile
func (c Service) RemoveProvisionWatcher(name string) (*common.BaseResponse, error) {
	client := c.getProvisionWatcherClient()

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	response, err := client.DeleteProvisionWatcherByName(context.Background(), name)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &response, nil
}

// UpdateDevice updates a device
func (c Service) UpdateDevice(name, id, description, service, profile, adminState, operState, location *string, labels []string, protocols map[string]dtos.ProtocolProperties) (*common.BaseResponse, error) {
	client := c.getDeviceClient()

	err := validateAdminState(*adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewUpdateDeviceRequest(dtos.UpdateDevice{
		Name:           name,
		Id:             id,
		Description:    description,
		ProfileName:    profile,
		ServiceName:    service,
		AdminState:     adminState,
		OperatingState: operState,
		Location:       location,
		Labels:         labels,
		Protocols:      protocols,
	})

	response, err := client.Update(context.Background(), []requests.UpdateDeviceRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Println(response[0])
	return &response[0], nil
}

func (c Service) UpdateProvisionWatcher(name, id, service, profile, adminState *string, labels []string, identifiers map[string]string) (*common.BaseResponse, error) {
	client := c.getProvisionWatcherClient()

	err := validateAdminState(*adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewUpdateProvisionWatcherRequest(dtos.UpdateProvisionWatcher{
		Name:        name,
		Id:          id,
		ServiceName: service,
		ProfileName: profile,
		AdminState:  adminState,
		Labels:      labels,
		Identifiers: identifiers,
	})

	response, err := client.Update(context.Background(), []requests.UpdateProvisionWatcherRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Println(response[0])
	return &response[0], nil

}

// UpdateDeviceService updates a device service
func (c Service) UpdateDeviceService(name, id, description, baseAddress, adminState *string, labels []string) (*common.BaseResponse, error) {
	client := c.getDeviceServiceClient()

	err := validateAdminState(*adminState)
	if err != nil {
		return nil, err
	}

	var req = requests.NewUpdateDeviceServiceRequest(dtos.UpdateDeviceService{
		Name:        name,
		Id:          id,
		Description: description,
		Labels:      labels,
		BaseAddress: baseAddress,
		AdminState:  adminState,
	})

	response, err := client.Update(context.Background(), []requests.UpdateDeviceServiceRequest{req})

	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Println(response[0])
	return &response[0], nil
}

func (c Service) getProvisionWatcherClient() interfaces.ProvisionWatcherClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewProvisionWatcherClient(url)
}

func (c Service) getDeviceClient() interfaces.DeviceClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceClient(url)
}

func (c Service) getDeviceServiceClient() interfaces.DeviceServiceClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceServiceClient(url)
}

func (c Service) getDeviceProfileClient() interfaces.DeviceProfileClient {
	url := fmt.Sprintf("http://%s:%v", c.Host, c.Port)
	return http.NewDeviceProfileClient(url)
}

func validateName(deviceService string) error {
	var err error
	if deviceService == "" {
		err = errors.New("name must be specified")
	}
	return err
}

func validateAdminState(adminState string) error {
	if !(adminState == models.Locked || adminState == models.Unlocked) {
		return fmt.Errorf("admin state should be %s or %s", models.Locked, models.Unlocked)
	}
	return nil
}

func validateOperatingState(operState string) error {
	if !(operState == models.Up || operState == models.Down || operState == models.Unknown) {
		return fmt.Errorf("operating state should be one of %s,%s or %s", models.Up, models.Down, models.Unknown)
	}
	return nil
}
