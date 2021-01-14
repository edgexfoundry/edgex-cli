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
	"github.com/BurntSushi/toml"
	"io"
	"os"
	"path/filepath"
)

const PathId = "/id/"
const PathName = "/name/"
const SampleConfigFileName = "sample-configuration.toml"
const ConfigFileName = "configuration.toml"

var DefaultConfigFile = filepath.Join(os.Getenv("HOME"), ".edgex-cli", ConfigFileName)
var Conf Configuration

// Configuration struct will use this to write config file eventually
type Configuration struct {
	Clients ClientInfo
}

type ClientInfo map[string]Client

type Client struct {
	Host     string
	Protocol string
	Port     int
}

func (c Client) Url() string {
	url := fmt.Sprintf("%s://%s:%v", c.Protocol, c.Host, c.Port)
	return url
}

// Security struct for security related config
type Security struct {
	Enabled bool
	Token   string
}

func LoadConfig(env Environment) error {
	var configFilePath string
	if env.IsSet("config-file") {
		configFilePath = env.GetString("config-file")
	} else {
		configFilePath = DefaultConfigFile
	}

	_, err := os.Stat(configFilePath)
	//Relative path differs depending on if the application is run from IDE or from distributed archive
	if os.IsNotExist(err) {
		_, err := copy("../res/"+SampleConfigFileName, configFilePath)
		if err != nil {
			_, err1 := copy("./res/"+SampleConfigFileName, configFilePath)
			if err1 != nil {
				fmt.Printf("failed to create configuration file '%s'. \n "+
					"%s\n %s\n ", configFilePath, err, err1)
				return errors.New(err.Error())
			}
		}
		fmt.Printf("Configuration file %s created\n", configFilePath)
	}

	if _, err := toml.DecodeFile(configFilePath, &Conf); err != nil {
		fmt.Printf("Error occurred while parsing %s: %s", configFilePath, err)
		return errors.New(err.Error())
	}
	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(dst), 0744)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
