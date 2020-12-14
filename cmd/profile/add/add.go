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
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"path/filepath"
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
    "DeviceResources": {{if .DeviceResources}}{{EscapeHTML .DeviceResources}}{{end}} {{if not .DeviceResources}}[]{{end}},
    "DeviceCommands":  {{if .DeviceCommands}}{{EscapeHTML .DeviceCommands}}{{end}} {{if not .DeviceCommands}}[]{{end}},
    "CoreCommands":  {{if .CoreCommands}}{{EscapeHTML .CoreCommands}}{{end}} {{if not .CoreCommands}}[]{{end}}
 }`

var interactiveMode bool
var file string

var name string
var description string
var manufacturer string
var model string
var labels string

// NewCommand returns the update device profile command
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
		return createProfilesFromFile(cmd)
	}

	profiles, err := parseProfile(interactiveMode)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	for _, p := range profiles {
		_, err = metadata.NewDeviceProfileClient(client).Add(cmd.Context(), &p)
		if err != nil {
			return err
		}
	}

	return nil
}

func createProfilesFromFile(cmd *cobra.Command) error {
	profiles, err := LoadFromFile(file)
	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute)
	for _, d := range profiles {
		_, err = metadata.NewDeviceProfileClient(client).Add(cmd.Context(), &d)
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
			"lastElem":   editor.IsLastElementOfSlice,
			"EscapeHTML": editor.EscapeHTML,
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

//LoadFromFile could read a file (json/yaml) that contains single DeviceProfile or list of DeviceProfile
func LoadFromFile(fPath string) ([]models.DeviceProfile, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json/Yaml")
		}
	}()
	file, err := ioutil.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	var profiles []models.DeviceProfile
	fExtension := strings.ToLower(filepath.Ext(fPath))
	err = errors.New("supported file format are yaml and json")
	if fExtension == ".yaml" || fExtension == ".yml" {
		err = unmarshalYaml(file, &profiles)
	}

	if fExtension == ".json" {
		err = unmarshalJSON(file, &profiles)
	}
	return profiles, err
}

//unmarshalJSON checks if the file contains just one Device Profile or list of Device Profiles in json format
func unmarshalJSON(file []byte, profiles *[]models.DeviceProfile) error {
	var p models.DeviceProfile
	var err error
	//check if the file contains a Device Profile in Json format
	if err = json.Unmarshal(file, &p); err != nil {
		//check if the file contains list of Device Profiles in Json format
		if err = json.Unmarshal(file, profiles); err != nil {
			return err
		}
	} else {
		*profiles = append(*profiles, p)
	}
	return nil
}

//unmarshalYaml checks if the file contains just one Device Profile or list of Device Profiles in yaml format
func unmarshalYaml(file []byte, profiles *[]models.DeviceProfile) error {
	var p models.DeviceProfile
	var err error
	//Then check if the file contains a Device Profile in Yaml format
	if err = yaml.Unmarshal(file, &p); err != nil {
		//Then check if the file contains list of Device Profiles in Yaml format
		if err = yaml.Unmarshal(file, profiles); err != nil {
			return err
		}
	} else {
		*profiles = append(*profiles, p)
	}
	return nil
}
