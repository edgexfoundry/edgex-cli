/*******************************************************************************
 * Copyright 2020 VMWare.
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/editor"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

const watcherTemplate = `[{{range $d := .}} 
{
  "Name": "{{.Name}}",
  "AdminState":"{{.AdminState}}",
  "OperatingState": "{{.OperatingState}}",
  "Service": {
     "Name": "{{.Service.Name}}"
  },
  "Profile": {
     "Name": "{{.Profile.Name}}"
  },
  "Identifiers":{ {{$s := separator ", "}}{{range $key, $value := .Identifiers}}{{call $s}}
     "{{$key}}": "{{$value}}"{{end}} 
  },
  "BlockingIdentifiers": { {{$s := separator ", "}}{{range $key, $value := .BlockingIdentifiers}}{{call $s}}
     "{{$key}}": [{{range $k, $v := $value}}"{{$v}}"{{end}}]{{end}} 
  }
}
{{end}}
]`

var interactiveMode bool
var name string
var adminState string
var profileName string
var serviceName string
var numberOfIdentifiers int
var numberOfBlockingIdentifiers int
var file string

// NewCommand returns the add watcher command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add watcher(s)",
		Long:  `Create watchers described in a given JSON file or use the interactive mode with additional flags.`,
		RunE:  handler,
	}
	cmd.Flags().BoolVarP(&interactiveMode, editor.InteractiveModeLabel, "i", false, "Open a default editor to customize the watcher information")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Watcher Name")
	cmd.Flags().StringVarP(&adminState, "adminState", "a", "", "Admin Status")
	cmd.Flags().StringVarP(&profileName, "profileName", "p", "", "Device Profile name")
	cmd.Flags().StringVarP(&serviceName, "serviceName", "s", "", "Device Service name")
	cmd.Flags().IntVar(&numberOfIdentifiers, "identifiers", 0, "Number of Identifiers")
	cmd.Flags().IntVarP(&numberOfBlockingIdentifiers, "blockingIdentifiers", "b", 0, "Number of BlockingIdentifiers")

	cmd.Flags().StringVarP(&file, "file", "f", "", "File containing watcher configuration in json format")
	return cmd
}

func handler(cmd *cobra.Command, args []string) (err error) {
	if interactiveMode && file != "" {
		return errors.New("you could work with interactive mode or file, but not with both")
	}

	var watchers []models.ProvisionWatcher
	if file != "" {
		watchers, err = LoadFromFile(file)
	} else {
		watchers, err = parse(interactiveMode)

	}

	if err != nil {
		return err
	}

	client := local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiProvisionWatcherRoute)
	for _, w := range watchers {
		_, err = metadata.NewProvisionWatcherClient(client).Add(context.Background(), &w)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	return nil
}

//LoadFromFile could read a file that contains single ProvisionWatcher or list of ProvisionWatcher
func LoadFromFile(filePath string) ([]models.ProvisionWatcher, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: Invalid Json")
		}
	}()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var watchers []models.ProvisionWatcher

	//check if the file contains just one ProvisionWatcher
	var w models.ProvisionWatcher
	err = json.Unmarshal(file, &w)
	if err != nil {
		//check if the file contains list of ProvisionWatcher
		err = json.Unmarshal(file, &watchers)
		if err != nil {
			return nil, err
		}
	} else {
		watchers = append(watchers, w)
	}
	return watchers, nil
}

func parse(interactiveMode bool) ([]models.ProvisionWatcher, error) {
	//parse ProvisionWatcher based on interactive mode and the other provided flags
	var err error
	var watchers []models.ProvisionWatcher
	populate(&watchers)
	var updatedWatchersBytes []byte
	if interactiveMode {
		updatedWatchersBytes, err = editor.OpenInteractiveEditor(watchers, watcherTemplate, template.FuncMap{
			"separator": separator,
		})
	}
	if err != nil {
		return nil, err
	}

	var updatedWatchers []models.ProvisionWatcher
	err = json.Unmarshal(updatedWatchersBytes, &updatedWatchers)
	if err != nil {
		return nil, errors.New("Unable to execute the command. The provided information is not valid:" + err.Error())
	}
	return updatedWatchers, err
}

func populate(watchers *[]models.ProvisionWatcher) {
	w := models.ProvisionWatcher{}
	w.Name = name
	w.AdminState = models.AdminState(adminState)
	w.Profile = models.DeviceProfile{Name: profileName}
	w.Service = models.DeviceService{Name: serviceName}
	if numberOfIdentifiers > 0 {
		identifier := make(map[string]string, numberOfIdentifiers)
		for i := 0; i < numberOfIdentifiers; i++ {
			identifier[fmt.Sprintf("Identifier-%d", i)] = "Example Value"
		}
		w.Identifiers = identifier
	}
	if numberOfBlockingIdentifiers > 0 {
		bi := make(map[string][]string, numberOfBlockingIdentifiers)
		for i := 0; i < numberOfBlockingIdentifiers; i++ {
			bi[fmt.Sprintf("Blocking Identifier-%d", i)] = []string{"Example Value-1"}
		}
		w.BlockingIdentifiers = bi
	}
	*watchers = append(*watchers, w)
}

func separator(s string) func() string {
	i := -1
	return func() string {
		i++
		if i == 0 {
			return ""
		}
		return s
	}
}
