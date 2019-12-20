/*******************************************************************************
 * Copyright 2019 VMware Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"

	"github.com/edgexfoundry-holding/edgex-cli/config/mocks"
)

var TestConfigFilePath = "config.yaml"
var TestConfigDirPath = "testdata/"
var TestCompletePath = TestConfigDirPath + TestConfigFilePath
var TestInvalidConfigFilePath = ""
var Error = errors.New("test error")

var TestConf = Configuration{
	Host: "localhost",
	SchedulerService: SchedulerService{
		Port:                          "48085",
		IntervalByIDRoute:             "interval",
		IntervalByNameSlugRoute:       "interval/name/",
		IntervalActionByIDRoute:       "intervalaction/",
		IntervalActionByNameSlugRoute: "intervalaction/name/",
	},
	NotificationService: NotificationService{
		Port:                        "48060",
		SubscriptionByIDRoute:       "subscription",
		SubscriptionByNameSlugRoute: "subscription/name/",
		NotificationByAgeRoute:      "notification/age/",
		NotificationByNameSlugRoute: "notification/slug/",
	},
	MetadataService: MetadataService{
		Port:                         "48081",
		DeviceServiceByIDRoute:       "deviceservice/id/",
		DeviceServiceBySlugNameRoute: "deviceservice/name/",
		DeviceByIDRoute:              "device/id/",
		DeviceBySlugNameRoute:        "device/name/",
		DeviceProfileByIDRoute:       "deviceprofile/id/",
		DeviceProfileBySlugNameRoute: "deviceprofile/name/",
		AddressableList:              "addressable",
	},
	DataService: DataService{
		Port:                       "48080",
		ReadingByIDRoute:           "reading/id/",
		VDescriptorByIDRoute:       "valuedescriptor/id/",
		VDescriptorByNameRoute:     "valuedescriptor/name/",
		DeleteEventByDeviceIDRoute: "event/device/",
	},
	ExportService: ExportService{
		Port:                    "48071",
		RegistrationByIDRoute:   "registration/",
		RegistrationByNameRoute: "registration/name/",
	},
}

func TestSetConfig(t *testing.T) {
	tests := []struct {
		name              string
		env               Environment
		configFilePath    string
		configDirPath     string
		expectError       bool
		expectedErrorType error
	}{
		{
			name:              "Error ReadInConfig",
			env:               createMockEnvReadInConfigErr(),
			configFilePath:    TestConfigFilePath,
			configDirPath:     TestConfigDirPath,
			expectError:       true,
			expectedErrorType: Error,
		},
		{
			name:              "Successful SetConfig",
			env:               createMockEnvSuccess(),
			configFilePath:    TestConfigFilePath,
			configDirPath:     TestConfigDirPath,
			expectError:       false,
			expectedErrorType: nil,
		},

		{
			name:              "Error Unmarshal",
			env:               createMockEnvUnmarshalErr(),
			configFilePath:    TestConfigFilePath,
			configDirPath:     TestConfigDirPath,
			expectError:       true,
			expectedErrorType: Error,
		},
		{
			name:              "Error WriteConfig",
			env:               createMockEnvWriteConfigErr(),
			configFilePath:    TestConfigFilePath,
			configDirPath:     TestConfigDirPath,
			expectError:       true,
			expectedErrorType: Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetConfig(test.env, test.configDirPath, test.configFilePath)

			if test.expectError && err == nil {
				t.Error("We expected an error but did not get one")
			}

			if !test.expectError && err != nil {
				t.Errorf("We do not expected an error but got one. %s", err.Error())
			}

			if test.expectError {
				eet := reflect.TypeOf(test.expectedErrorType)
				aet := reflect.TypeOf(err)
				if !aet.AssignableTo(eet) {
					t.Errorf("Expected error of type %v, but got an error of type %v", eet, aet)
				}
			}

			return
		})
	}
}

func createMockEnvSuccess() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", TestCompletePath).Return(nil)
	dbMock.On("ReadInConfig").Return(nil)
	dbMock.On("Unmarshal", &TestConf).Return(nil)
	dbMock.On("WriteConfig").Return(nil)
	return &dbMock
}

func createMockEnvReadInConfigErr() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", TestCompletePath).Return(nil)
	dbMock.On("ReadInConfig").Return(Error)
	dbMock.On("Unmarshal", &TestConf).Return(nil)
	dbMock.On("WriteConfig").Return(nil)
	return &dbMock
}

func createMockEnvUnmarshalErr() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", TestCompletePath).Return(nil)
	dbMock.On("ReadInConfig").Return(nil)
	dbMock.On("Unmarshal", &TestConf).Return(Error)
	dbMock.On("WriteConfig").Return(nil)
	return &dbMock
}

func createMockEnvWriteConfigErr() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", TestCompletePath).Return(nil)
	dbMock.On("ReadInConfig").Return(nil)
	dbMock.On("Unmarshal", &TestConf).Return(nil)
	dbMock.On("WriteConfig").Return(Error)
	return &dbMock
}

func createFsMockSuccess() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Write", mock.Anything).Return(0, nil)
	fileMock.On("Close").Return(nil)
	fsMock.On("Create", TestConfigFilePath).Return(&fileMock, nil)
	return &fsMock
}

func createFsMockErr() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Write", mock.Anything).Return(0, nil)
	fileMock.On("Close").Return(nil)
	fsMock.On("Create", TestConfigFilePath).Return(&fileMock, Error)
	return &fsMock
}

func createFsMockFileCloseErr() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Write", mock.Anything).Return(0, nil)
	fileMock.On("Close").Return(Error)
	fsMock.On("Create", TestConfigFilePath).Return(&fileMock, nil)
	return &fsMock
}

func createFsMockFileWriteErr() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Write", mock.Anything).Return(0, Error)
	fileMock.On("Close").Return(Error)
	fsMock.On("Create", TestConfigFilePath).Return(&fileMock, nil)
	return &fsMock
}

func TestCreateDefaultFile(t *testing.T) {
	tests := []struct {
		name              string
		fsMock            afero.Fs
		configFilePath    string
		configuration     *Configuration
		expectError       bool
		expectedErrorType error
	}{
		{
			name:              "Successful creation",
			fsMock:            createFsMockSuccess(),
			configFilePath:    TestConfigFilePath,
			configuration:     &TestConf,
			expectError:       false,
			expectedErrorType: nil,
		},
		{
			name:              "Error Create",
			fsMock:            createFsMockErr(),
			configFilePath:    TestConfigFilePath,
			configuration:     &TestConf,
			expectError:       true,
			expectedErrorType: Error,
		},
		{
			name:              "Error File Write",
			fsMock:            createFsMockFileWriteErr(),
			configFilePath:    TestConfigFilePath,
			configuration:     &TestConf,
			expectError:       true,
			expectedErrorType: Error,
		},
		{
			name:              "Error File Close",
			fsMock:            createFsMockFileCloseErr(),
			configFilePath:    TestConfigFilePath,
			configuration:     &TestConf,
			expectError:       true,
			expectedErrorType: Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := createDefaultFile(test.configFilePath, test.configuration, test.fsMock)

			if test.expectError && err == nil {
				t.Error("We expected an error but did not get one")
			}

			if !test.expectError && err != nil {
				t.Errorf("We do not expected an error but got one. %s", err.Error())
			}

			if test.expectError {
				eet := reflect.TypeOf(test.expectedErrorType)
				aet := reflect.TypeOf(err)
				if !aet.AssignableTo(eet) {
					t.Errorf("Expected error of type %v, but got an error of type %v", eet, aet)
				}
			}

			return
		})
	}
}
