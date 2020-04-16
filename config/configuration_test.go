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
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/edgexfoundry-holding/edgex-cli/config/mocks"
)

var workDir, _ = os.Getwd()
var TestValidConfigFilePath = filepath.Join(workDir,  "..", "res", "configuration.toml")
var TestInvalidConfigFilePath = filepath.Join(workDir, "..", "res", "invalidConfig.toml")
var Error = errors.New("test error")

var TestConf = Configuration{
	Clients: ClientInfo{
		"Clients.Metadata": Client{Host: "localhost", Protocol: "http", Port: 48081},
		"Clients.CoreData": Client{Host: "localhost", Protocol: "http", Port: 48080},
		"Clients.Scheduler": Client{Host: "localhost", Protocol: "http", Port: 48085},
		"Clients.Notification": Client{Host: "localhost", Protocol: "http", Port: 48060},
		"Clients.Logging": Client{Host: "localhost", Protocol: "http", Port: 48061},
	},
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name              string
		env               Environment
		configFilePath    string
		expectError       bool
		expectedErrorType error
	}{
		{
			name:              "Successful GetDefaultConfig",
			env:               getDefaultConfigFileMockEnvSuccess(),
			configFilePath:    TestValidConfigFilePath,
			expectError:       false,
			expectedErrorType: nil,
		},
		{
			name:              "Successful GetConfig",
			env:               getConfigFileMockEnvSuccess(),
			configFilePath:    TestValidConfigFilePath,
			expectError:       false,
			expectedErrorType: nil,
		},
/*		{
			name:              "Unsuccessful GetConfig",
			env:               getConfigFileMockEnvError(),
			configFilePath:    TestInvalidConfigFilePath,
			expectError:       true,
			expectedErrorType: nil,
		},*/
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err = LoadConfig(test.env)

			if test.expectError && err == nil {
				t.Error("We expected an error but did not get one")
			}

			if !test.expectError && err != nil {
				t.Errorf("We did not expect an error but got one. %s", err.Error())
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


func getDefaultConfigFileMockEnvSuccess() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", DefaultConfigFile).Return(nil)
	dbMock.On("GetString").Return(DefaultConfigFile)
	dbMock.On("IsSet", mock.Anything).Return(false)
	return &dbMock
}
func getConfigFileMockEnvSuccess() Environment {
	dbMock := mocks.Environment{}

	dbMock.On("SetConfigFile", TestValidConfigFilePath).Return(nil)
	dbMock.On("GetString", mock.Anything).Return(TestValidConfigFilePath)
	dbMock.On("IsSet", mock.Anything).Return(true)
	return &dbMock
}

func getConfigFileMockEnvError() Environment {
	dbMock := mocks.Environment{}

	dbMock.On("SetConfigFile", TestInvalidConfigFilePath).Return(nil)
	dbMock.On("GetString", mock.Anything).Return(TestInvalidConfigFilePath)
	dbMock.On("IsSet",mock.Anything).Return(true)
	return &dbMock
}

/*
func readFsMockSuccess() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Read", TestValidConfigFilePath).Return(&fileMock, nil)
	fileMock.On("Close").Return(nil)
	return &fsMock
}

func readFsMockNonExistentFile() afero.Fs {
	fsMock := mocks.Fs{}
	fileMock := mocks.File{}
	fileMock.On("Read", TestInvalidConfigFilePath).Return(nil, Error)
	fileMock.On("Close").Return(nil)
	return &fsMock
}

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name              string
		fsMock            afero.Fs
		env               Environment
		configFilePath    string
		expectError       bool
		expectedErrorType error
	}{
		{
			name:              "Successful read",
			env:                Environment,
			fsMock:            readFsMockSuccess(),
			configFilePath:    TestValidConfigFilePath,
			expectError:       false,
			expectedErrorType: nil,
		},
		{
			name:              "Error read",
			env:                Environment,
			fsMock:            readFsMockErr(),
			configFilePath:    TestInvalidConfigFilePath,
			expectError:       true,
			expectedErrorType: Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LoadConfig(test.env)

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
}*/
