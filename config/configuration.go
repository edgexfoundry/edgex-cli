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
	"fmt"
	"os"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var AppFs = afero.NewOsFs()

// Configuration struct will use this to write config file eventually
type Configuration struct {
	Security
	//Client 				map[string]ClientInfo
	DataService         DataService
	MetadataService     MetadataService
	SchedulerService    SchedulerService
	ExportService       ExportService
	NotificationService NotificationService
}

// Security struct for security related config
type Security struct {
	Enabled bool
	Token   string
}

type SchedulerService struct {
	Port     string
	Host     string
	Protocol string
}

type NotificationService struct {
	Port     string
	Host     string
	Protocol string
}

type MetadataService struct {
	Port     string
	Host     string
	Protocol string
}

type DataService struct {
	Port     string
	Host     string
	Protocol string
}

type ExportService struct {
	Port     string
	Host     string
	Protocol string
}

var Conf Configuration = Configuration{
	SchedulerService: SchedulerService{
		Port:     "48085",
		Host:     "localhost",
		Protocol: "http",
	},
	NotificationService: NotificationService{
		Port:     "48060",
		Host:     "localhost",
		Protocol: "http",
	},
	MetadataService: MetadataService{
		Port:     "48081",
		Host:     "localhost",
		Protocol: "http",
	},
	DataService: DataService{
		Port:     "48080",
		Host:     "localhost",
		Protocol: "http",
	},
	ExportService: ExportService{
		Port:     "48071",
		Host:     "localhost",
		Protocol: "http",
	},
}

// SetConfig create the config file if doesn't exists
func SetConfig(env Environment, configDirPath string, configFilePath string) error {

	configuration := &Conf
	completePath := configDirPath + configFilePath

	if !exists(completePath) {
		if !exists(configDirPath) {
			err := os.Mkdir(configDirPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		err := createDefaultFile(completePath, configuration, AppFs)
		if err != nil {
			return err
		}
	}

	env.SetConfigFile(completePath)

	// Reading from file that was already existing or newly created
	if err := env.ReadInConfig(); err != nil && exists(completePath) {
		return fmt.Errorf("error reading config file, %s", err)
	}

	err := env.Unmarshal(configuration)
	if err != nil {
		fmt.Errorf("unable to decode into struct")
		return err
	}

	err = env.WriteConfig()
	if err != nil {
		fmt.Errorf("unable to write to viper config")
		return err
	}

	return nil
}

// Helper function to check whether file exists
func exists(configPath string) bool {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createDefaultFile(configFilePath string, configuration *Configuration, fs afero.Fs) error {
	y, err := yaml.Marshal(configuration)
	if err != nil {
		return err
	}

	f, err := fs.Create(configFilePath)
	if err != nil {
		return err
	}
	_, err = f.Write(y)
	if err != nil {
		err = f.Close()
		if err != nil {
			return err
		}

	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
