// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"encoding/json"
	"fmt"
	"time"
	"bytes"
	"io/ioutil"
	"net/http"
	"errors"

	models "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
	"github.com/pelletier/go-toml"
)

type DeviceConfig struct {
	Name string
	Profile string
	Description string
	Labels []string
	Protocols map[string]models.ProtocolProperties
	AutoEvents []models.AutoEvent
}

type Profile struct {
	Name string
}

type Service struct {
	Name string
}

type Device struct {
	Name string
	Origin int64
	Profile Profile
	Service Service
	Description string
	Labels []string
	AdminState models.AdminState
	OperatingState models.OperatingState
	Protocols map[string]models.ProtocolProperties
	AutoEvents []models.AutoEvent
}

type DeviceFile struct {
	Service string
	DeviceList []DeviceConfig
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add devices",
		Long:  `Create the devices described in the given TOML files.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add Devices:")
			for _, val := range args {
				fmt.Println(val, "... ")
				processFile(val)
			}
		},
	}
	return cmd
}

func addDevice (dev Device) (string, error) {
	jsonStr, err := json.Marshal (dev)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:48081/api/v1/device", bytes.NewBuffer (jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return string(respBody), nil
	} else {
		return "", errors.New (string(respBody))
	}
}

func processFile (fname string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("...Invalid TOML")
		}
	}()

	var content = &DeviceFile{}
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("...Error loading file: ", err.Error())
		return
	}

	err = toml.Unmarshal(file, content)
	if err != nil {
		fmt.Println("...Error parsing file: ", err.Error())
		return
	}

	for _, d := range content.DeviceList {
		fmt.Println("...Create device ", d.Name)
		prf := Profile {
			Name:	d.Profile,
		}
		svc := Service {
			Name:	content.Service,
		}
		millis := time.Now().UnixNano() / int64(time.Millisecond)
		dev := Device {
			Name:		d.Name,
			Profile:	prf,
			Protocols:	d.Protocols,
			Labels:		d.Labels,
			Service:	svc,
			AdminState:	models.Unlocked,
			OperatingState:	models.Enabled,
			AutoEvents:	d.AutoEvents,
		}
		dev.Origin = millis
		id, err := addDevice (dev)
		if err != nil {
			fmt.Println("......Error: ", err.Error())
		} else {
			fmt.Println("......Created with ID ", id)
		}
	}
}
