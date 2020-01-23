// Copyright © 2019 VMware, INC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package add

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type NotificationsFile struct {
	Notifications []models.Notification
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add notification",
		Long:  `Create the notification described in the given TOML files.`,
		Run:   addNotificationHandler,
	}
	return cmd
}

func addNotificationHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Add Notification:")
	for _, val := range args {
		fmt.Println(val, "... ")
		processFile(val)
	}
}

func addNotification(n *models.Notification) (string, error) {
	jsonStr, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://" + viper.GetString("Host") + ":48060/api/v1/notification", bytes.NewBuffer(jsonStr))
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

	if resp.StatusCode == 202 {
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

	var content = &NotificationsFile{}
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
	for _, n := range content.Notifications {
		fmt.Println("...Create notification ", n)
		id, err := addNotification(&n)
		if err != nil {
			fmt.Println("......Error: ", err.Error())
		} else {
			fmt.Println("......Created with ID ", id)
		}
	}
}
