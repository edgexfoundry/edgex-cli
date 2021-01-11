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
	"fmt"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/edgexfoundry/edgex-cli/config/mocks"
)

var workDir, _ = os.Getwd()
var ValidConfigFile = filepath.Join(workDir, "..", "res", "configuration.toml")
var NonExistentConfigFile = filepath.Join(workDir, "testdata", "nonExistentConfig.toml")
var InvalidTomlConfigFile = filepath.Join(workDir, "testdata", "invalidConfig.toml")
var Error = errors.New("test error")

var TestConf = Configuration{
	Clients: ClientInfo{
		"Clients.Metadata":     Client{Host: "localhost", Protocol: "http", Port: 48081},
		"Clients.CoreData":     Client{Host: "localhost", Protocol: "http", Port: 48080},
		"Clients.Scheduler":    Client{Host: "localhost", Protocol: "http", Port: 48085},
		"Clients.Notification": Client{Host: "localhost", Protocol: "http", Port: 48060},
		"Clients.Logging":      Client{Host: "localhost", Protocol: "http", Port: 48061},
	},
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name           string
		env            Environment
		configFilePath string
		//result            Configuration
		expectError       bool
		expectedErrorType error
	}{
		{
			name:              "Successful GetDefaultConfig",
			env:               getDefaultConfigFileMockEnvSuccess(),
			configFilePath:    DefaultConfigFile,
			expectError:       false,
			expectedErrorType: nil,
		},
		{
			name:              "Successful GetConfig",
			env:               getConfigFileMockEnvSuccess(),
			configFilePath:    ValidConfigFile,
			expectError:       false,
			expectedErrorType: nil,
		},
		{
			name:              "Unsuccessful decode Config",
			env:               getConfigFileMockDecodeError(),
			configFilePath:    InvalidTomlConfigFile,
			expectError:       true,
			expectedErrorType: Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err = LoadConfig(test.env)
			fmt.Println(err)
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
		})
	}
}

func getDefaultConfigFileMockEnvSuccess() Environment {
	dbMock := mocks.Environment{}
	dbMock.On("SetConfigFile", DefaultConfigFile).Return(nil)
	dbMock.On("GetString", "config-file").Return(DefaultConfigFile)
	dbMock.On("IsSet", mock.Anything).Return(false)
	return &dbMock
}
func getConfigFileMockEnvSuccess() Environment {
	dbMock := mocks.Environment{}

	dbMock.On("SetConfigFile", ValidConfigFile).Return(nil)
	dbMock.On("GetString", "config-file").Return(ValidConfigFile)
	dbMock.On("IsSet", mock.Anything).Return(true)
	return &dbMock
}

func getConfigFileMockEnvError() Environment {
	dbMock := mocks.Environment{}

	dbMock.On("SetConfigFile", NonExistentConfigFile).Return(nil)
	dbMock.On("GetString", "config-file").Return(NonExistentConfigFile)
	dbMock.On("IsSet", mock.Anything).Return(true)
	return &dbMock
}

func getConfigFileMockDecodeError() Environment {
	dbMock := mocks.Environment{}

	dbMock.On("SetConfigFile", InvalidTomlConfigFile).Return(nil)
	dbMock.On("GetString", "config-file").Return(InvalidTomlConfigFile)
	dbMock.On("IsSet", mock.Anything).Return(true)
	return &dbMock
}
