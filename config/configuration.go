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

	"github.com/ghodss/yaml"
	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

// Configuration struct will use this to write config file eventually
type Configuration struct {
	Host string
	Security
	DataService      DataService
	MetadataService  MetadataService
	SchedulerService SchedulerService
	ExportService    ExportService
}

// Security struct for security related config
type Security struct {
	Enabled bool
	Token   string
}

type SchedulerService struct {
	Port                          string
	IntervalByIDRoute             string
	IntervalByNameSlugRoute       string
	IntervalActionByIDRoute       string
	IntervalActionByNameSlugRoute string
}

type MetadataService struct {
	Port                         string
	DeviceServiceByIDRoute       string
	DeviceServiceBySlugNameRoute string
	DeviceByIDRoute              string
	DeviceBySlugNameRoute        string
	DeviceProfileByIDRoute       string
	DeviceProfileBySlugNameRoute string
}

type DataService struct {
	Port                   string
	ReadingByIDRoute       string
	VDescriptorByIDRoute   string
	VDescriptorByNameRoute string
}

type ExportService struct {
	Port                    string
	RegistrationByIDRoute   string
	RegistrationByNameRoute string
}

var Conf Configuration = Configuration{
	SchedulerService: SchedulerService{
		Port:                          "48085",
		IntervalByIDRoute:             "interval",
		IntervalByNameSlugRoute:       "interval/name/",
		IntervalActionByIDRoute:       "intervalaction/",
		IntervalActionByNameSlugRoute: "intervalaction/name/",
	},
	MetadataService: MetadataService{
		Port:                         "48081",
		DeviceServiceByIDRoute:       "deviceservice/id/",
		DeviceServiceBySlugNameRoute: "deviceservice/name/",
		DeviceByIDRoute:              "device/id/",
		DeviceBySlugNameRoute:        "device/name/",
		DeviceProfileByIDRoute:       "deviceprofile/id/",
		DeviceProfileBySlugNameRoute: "deviceprofile/name/",
	},
	DataService: DataService{
		Port:                   "48080",
		ReadingByIDRoute:       "reading/id/",
		VDescriptorByIDRoute:   "valuedescriptor/id/",
		VDescriptorByNameRoute: "valuedescriptor/name/",
	},
	ExportService: ExportService{
		Port:                    "48071",
		RegistrationByIDRoute:   "registration/",
		RegistrationByNameRoute: "registration/name/",
	},
}

func SetConfig(env Environment, configFilePath string) error {

	var configuration *Configuration = &Conf
	if !exists(configFilePath) {

		err := createDefaultFile(configFilePath, configuration, AppFs)
		if err != nil {
			return err
		}
	}

	env.SetConfigFile(configFilePath)

	// Reading from file that was already existing or newly created
	if err := env.ReadInConfig(); err != nil && exists(configFilePath) {
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
		return  err
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
