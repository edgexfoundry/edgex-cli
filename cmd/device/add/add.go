// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
)

type DeviceConfig struct {
	Name            string
	Profile         string
	Description     string
	Service         string
	Labels          []string
	AddressableName string
	Protocols       map[string]models.ProtocolProperties
	AutoEvents      []models.AutoEvent
}

type DeviceFile struct {
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

func addDevice(dev models.Device) (string, error) {
	jsonStr, err := json.Marshal(dev)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonStr))
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:48081/api/v1/device", bytes.NewBuffer(jsonStr))
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
		return "", errors.New(string(respBody))
	}
}

func processFile(fname string) {
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
		millis := time.Now().UnixNano() / int64(time.Millisecond)
		dev := models.Device{
			Name:           d.Name,
			Profile:        models.DeviceProfile{Name: d.Profile},
			Protocols:      d.Protocols,
			Labels:         d.Labels,
			Service:        models.DeviceService{Name: d.Service, Addressable: models.Addressable{Name: d.AddressableName}},
			AdminState:     models.Unlocked,
			OperatingState: models.Enabled,
			AutoEvents:     d.AutoEvents,
		}
		dev.Origin = millis
		id, err := addDevice(dev)
		if err != nil {
			fmt.Println("......Error: ", err.Error())
		} else {
			fmt.Println("......Created with ID ", id)
		}
	}
}
