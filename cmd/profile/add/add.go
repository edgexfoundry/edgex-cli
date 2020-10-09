// Copyright © 2020 VMware, INC
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

const ProfilesTemplate = `[{{range $d := .}}` + ProfileTemplate + `{{end}}]`

const ProfileTemplate = `{
   "Description": "{{.Description}}",
   "Name" : "{{.Name}}",
   "Manufacturer" : "{{.Manufacturer}}",
   "Model" : "{{.Model}}",
   "Labels" : [
    {{- $labelsLenght := len .Labels}}
    {{- range $idx, $l := .Labels}}
      "{{$l}}" {{if not (lastElem $idx $labelsLenght)}},{{end}} 
    {{- end}}
    ],
    "DeviceResources": {{.DeviceResources}},
    "DeviceCommands": {{.DeviceCommands}},
    "CoreCommands": {{.CoreCommands}}
 }`

var interactiveMode bool
var file string

var name string
var description string
var manufacturer string
var model string
var labels string

// NewCommand returns the update deviceprofile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add profiles",
		Long:  `Create device profiles described in a given JSON file or use the interactive mode with additional flags.`,
		RunE:  handler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the Event information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVar(&manufacturer, "manufacturer", "", "Manufacturer")
	cmd.Flags().StringVar(&model, "model", "", "Model")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma separated strings")

	cmd.Flags().StringVarP(&file, "file", "f", "", "File containing profile configuration in json format")
	return cmd
}

func handler(cmd *cobra.Command, args []string) error {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	if file != "" {
		return createProfilesFromFile()
	}

	profiles, err := parseProfile(interactiveMode)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	for _, p := range profiles {
		_, err = metadata.NewDeviceProfileClient(client).Add(context.Background(), &p)
		if err != nil {
			return err
		}
	}

	return nil
}

func createProfilesFromFile() error {
	profiles, err := LoadProfilesFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	for _, d := range profiles {
		_, err = metadata.NewDeviceProfileClient(client).Add(context.Background(), &d)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	return nil
}

func parseProfile(interactiveMode bool) ([]models.DeviceProfile, error) {
	//parse DeviceProfile based on interactive mode and the other provided flags
	var err error
	var profiles []models.DeviceProfile
	populateProfile(&profiles)

	var updatedProfilesBytes []byte
	if interactiveMode {
		updatedProfilesBytes, err = editor.OpenInteractiveEditor(profiles, ProfilesTemplate, template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
			"lastElem": editor.IsLastElementOfSlice,
		})
	}
	if err != nil {
		return nil, err
	}

	var updatedProfiles []models.DeviceProfile
	err = json.Unmarshal(updatedProfilesBytes, &updatedProfiles)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is not valid:" + err.Error())
	}
	return updatedProfiles, err
}

func populateProfile(profiles *[]models.DeviceProfile) {
	d := models.DeviceProfile{}
	d.Name = name
	d.Description = description
	d.Model = model
	d.Manufacturer = manufacturer
	ls := strings.Split(labels, ",")
	for i := range ls {
		ls[i] = strings.TrimSpace(ls[i])
	}
	d.Labels = ls
	*profiles = append(*profiles, d)
}

//LoadProfilesFromFile could read a file that contains single DeviceProfile or list of DeviceProfile
func LoadProfilesFromFile(filePath string) ([]models.DeviceProfile, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var profiles []models.DeviceProfile

	//check if the file contains just one DeviceProfile
	var p models.DeviceProfile
	err = json.Unmarshal(file, &p)
	if err != nil {
		//check if the file contains list of DeviceProfile
		err = json.Unmarshal(file, &profiles)
		if err != nil {
			return nil, err
		}
	} else {
		profiles = append(profiles, p)
	}
	return profiles, nil
}
