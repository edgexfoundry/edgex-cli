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
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

const PathId ="/id/"
const PathName ="/name/"

var confDir = ".edgex-cli"
var fileName = "configuration.toml"
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

func LoadConfig() {
	//TODO make it work for Windows
	configFile := strings.Join([]string{os.Getenv("HOME"), confDir, fileName}, "/")
	if len(viper.GetString("config-file")) > 0 {
		configFile = viper.GetString("config-file")
	}

	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal("Error occured while parsing %s:%s", configFile, err)
	}
}
