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
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var categories = []string{models.Security, models.Hwhealth, models.Swhealth}
var severities = []string{models.Normal, models.Critical}
var statuses = []string{models.New, models.Processed, models.Escalated}

const NotificationsTemplate = `[{{range $d := .}}` + NotificationTemplate + `{{end}}]`
const NotificationTemplate = `{
   "Slug": "{{.Slug}}",
   "Sender" : "{{.Sender}}",
   "Category" : "{{.Category}}",
   "Severity" : "{{.Severity}}",
   "Content" : "{{.Content}}",
   "Description" : "{{.Description}}",
   "ContentType" : "{{.ContentType}}",
   "Status" : "{{.Status}}",
   "Labels" : [
    {{- $labelsLenght := len .Labels}}
    {{- range $idx, $l := .Labels}}
      "{{$l}}" {{if not (lastElem $idx $labelsLenght)}},{{end}} 
    {{- end}}
    ]
 }`

var interactiveMode bool
var slug string
var sender string
var category string
var severity string
var content string
var description string
var labels string
var contentType string
var status string

var file string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add notification",
		Long:  `Create notifications described in a given JSON file or use the interactive mode with additional flags.`,
		RunE:  addNotificationHandler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Notification information")
	cmd.Flags().StringVarP(&slug, "slug", "s", "", "Slug")
	cmd.Flags().StringVar(&sender, "sender", "", "Sender")
	cmd.Flags().StringVarP(&category, "category", "c", "", fmt.Sprintf("Status\nValues: [%s]", strings.Join(categories, ",")))
	cmd.Flags().StringVar(&severity, "severity", "", fmt.Sprintf("Status\nValues: [%s]", strings.Join(severities, ",")))
	cmd.Flags().StringVar(&content, "content", "", "Content")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVarP(&labels, "labels", "l", "", "Comma separated string")
	cmd.Flags().StringVar(&contentType, "contentType", "", "ContentType")
	cmd.Flags().StringVar(&status, "status", "", fmt.Sprintf("Status\nValues: [%s]", strings.Join(statuses, ",")))

	cmd.Flags().StringVarP(&file, "file", "f", "", "File containing notification configuration in json format")
	return cmd
}

func addNotificationHandler(cmd *cobra.Command, args []string) error {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	if file != "" {
		return createNotificationsFromFile()
	}

	notifications, err := parseNotification(interactiveMode)
	if err != nil {
		return err
	}

	createNotification(notifications)
	return nil
}

func parseNotification(interactiveMode bool) ([]models.Notification, error) {
	//parse Notification based on interactive mode and the other provided flags
	var err error
	var notifications []models.Notification
	populateNotification(&notifications)

	var updatedNotificationBytes []byte
	if interactiveMode {
		updatedNotificationBytes, err = editor.OpenInteractiveEditor(notifications, NotificationsTemplate, template.FuncMap{
			"lastElem": editor.IsLastElementOfSlice,
		})
	} else {
		updatedNotificationBytes, err = json.Marshal(notifications)
	}
	if err != nil {
		return nil, err
	}
	var updatedNotification []models.Notification
	err = json.Unmarshal(updatedNotificationBytes, &updatedNotification)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is not valid: " + err.Error())
	}
	return updatedNotification, err
}

func populateNotification(notifications *[]models.Notification) {
	i := models.Notification{}
	i.Slug = slug
	i.Sender = sender
	i.Category = models.NotificationsCategory(strings.ToUpper(category))
	i.Status = models.NotificationsStatus(strings.ToUpper(status))
	i.Severity = models.NotificationsSeverity(strings.ToUpper(severity))
	i.Content = content
	i.Description = description
	ls := strings.Split(labels, ",")
	for i := range ls {
		ls[i] = strings.TrimSpace(ls[i])
	}
	i.Labels = ls
	i.ContentType = contentType
	*notifications = append(*notifications, i)
}

func createNotificationsFromFile() error {
	notifications, err := LoadNotificationFromFile(file)
	if err != nil {
		return err
	}

	createNotification(notifications)
	return nil
}

//LoadNotificationFromFile could read a file that contains single Notification or list of Notifications
func LoadNotificationFromFile(filePath string) ([]models.Notification, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var notifications []models.Notification

	//check if the file contains just one Notification
	var n models.Notification
	err = json.Unmarshal(file, &n)
	if err != nil {
		//check if the file contains list of Notifications
		err = json.Unmarshal(file, &notifications)
		if err != nil {
			return nil, err
		}
	} else {
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func createNotification(notifications []models.Notification) {
	url := config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute
	for _, n := range notifications {
		request.Post(url, n)
	}
}
