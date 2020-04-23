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
	"os"
	"path/filepath"
)

const PathId ="/id/"
const PathName ="/name/"
var  DefaultConfigFile = filepath.Join(os.Getenv("HOME"), ".edgex-cli", "configuration.toml")
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
	if _, err := toml.DecodeFile(configFilePath, &Conf); err != nil {
		fmt.Printf("Error occurred while parsing %s: %s", configFilePath, err)
		return errors.New(err.Error())
	}
	return nil
}

