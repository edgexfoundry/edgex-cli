/*******************************************************************************
 * Copyright 2020 Dell Inc.
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

package add

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/coredata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/spf13/cobra"
)

// EventTemplate is a Go template which is used to generate a JSON representation of an Event. The data associated with
// this template is a modes.Event
//
// This template requires a function 'lastElem' which accepts and index(int) and a slice of Readings and determines if
// the index is the last one in the array.
const EventTemplate = `
{
  "Device" : "{{.Device}}",
  "Created" : {{.Created}},
  "Modified" : {{.Modified}},
  "Origin" : {{.Origin}},
  "Pushed" : {{.Pushed}}{{if .Readings}},{{$readings := .Readings}}
  "Readings" : [{{range $index, $reading := .Readings}}
    {
      "Name" : "{{$reading.Name}}",
      "Device": "{{.Device}}",
      "Value": "{{$reading.Value}}",
      "ValueType": "",
      "FloatEncoding": "",
      "BinaryValue": "",
      "MediaType": "",
      "Origin": {{.Origin}},
      "Pushed": {{.Pushed}},
      "Created": {{.Created}},
      "Origin": {{.Origin}},
      "Modified": {{.Modified}}
    } {{if not (lastElem $index $readings)}},{{end}}{{end}}
  ]{{end}}
}`

// Global variables that hold information passed in as flags.
var pushed int64
var origin int64
var created int64
var modified int64
var device string
var numberOfReadings int
var interactiveMode bool

// NewCommand returns the add event command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{

		Use:   "add",
		Short: "Create an event",
		Long:  `Create a new event`,
		RunE: func(cmd *cobra.Command, args []string) error {

			url := config.Conf.Clients["CoreData"].Url() + clients.ApiEventRoute
			newEvent := models.Event{}
			populateEvent(&newEvent)
			var newEventBytes []byte
			var err error
			if interactiveMode {
				newEventBytes, err = openInteractiveEditor(newEvent)
			} else {
				newEventBytes, err = json.Marshal(newEvent)
			}

			if err != nil {
				return err
			}

			updatedEvent := models.Event{}
			err = json.Unmarshal(newEventBytes, &updatedEvent)
			if err != nil {
				return errors.New("Unable to create an Event. The provided information is not valid")
			}

			client := local.New(url)
			ec := coredata.NewEventClient(client)
			_, err = ec.Add(cmd.Context(), &newEvent)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Event information")
	cmd.Flags().StringVarP(&device, "device", "d", "", "Device to which the Event is associated")
	cmd.Flags().Int64VarP(&pushed, "pushed", "p", 0, "Mark the Event as Pushed")
	cmd.Flags().Int64VarP(&created, "created", "c", 0, "Created timestamp")
	cmd.Flags().Int64VarP(&modified, "modified", "m", 0, "Modified timestamp")
	cmd.Flags().Int64VarP(&origin, "origin", "o", 0, "Origin")
	cmd.Flags().IntVarP(&numberOfReadings, "readings", "r", 0, "Number of Readings for the Event")
	return cmd
}

// populateEvent adds data to the event that have been passed as CLI flags.
func populateEvent(event *models.Event) {
	event.Device = device
	event.Pushed = pushed
	event.Created = created
	event.Modified = modified
	event.Origin = origin
	if numberOfReadings > 0 {
		readings := make([]models.Reading, numberOfReadings)
		for i := 0; i < numberOfReadings; i++ {
			tempReading := models.Reading{
				Name:  fmt.Sprintf("Reading-%d", i),
				Value: "Example Value", // This is a required field
			}
			readings[i] = tempReading
		}
		event.Readings = readings
	}
}

// openInteractiveEditor opens the users default editor and with a JSON representation of the Event.
func openInteractiveEditor(event models.Event) ([]byte, error) {
	funcMap := template.FuncMap{
		"lastElem": isLastElementOfSlice,
	}

	eventJsonTemplate, err := template.New("Event").Funcs(funcMap).Parse(EventTemplate)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	err = eventJsonTemplate.Execute(buff, event)
	if err != nil {
		return nil, err
	}

	return editor.CaptureInputFromEditor(buff.Bytes())
}

// isLastElementOfSlice is a function which is used in the EventTemplate to determine the last element in a slice of
// Readings.
func isLastElementOfSlice(index int, arr []models.Reading) bool {
	return len(arr)-1 == index
}
