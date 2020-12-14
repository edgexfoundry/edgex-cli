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
package list

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/edgexfoundry/edgex-cli/config"
	"github.com/edgexfoundry/edgex-cli/pkg/formatters"

	"github.com/spf13/cobra"
)

const template = "ID\tName\tProfile\tService\tAdminState\n" +
	"{{range .}}" +
	"{{.Id}}\t{{.Name}}\t{{.Profile.Name}}\t{{.Service.Name}}\t{{.AdminState}}\n" +
	"{{end}}"

// NewCommand returns the list watcher command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of watchers",
		Long:  `Return a list of watchers or retrieve a watcher by id`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  listHandler,
	}
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	client := metadata.NewProvisionWatcherClient(
		local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiProvisionWatcherRoute),
	)

	var watchers []models.ProvisionWatcher
	if len(args) == 0 {
		watchers, err = client.ProvisionWatchers(cmd.Context())
	} else {
		watchers, err = provisionWatchers(cmd, client, args[0])
	}
	if err != nil {
		return err
	}

	formatter := formatters.NewFormatter(template, nil)
	err = formatter.Write(watchers)
	return
}

func provisionWatchers(cmd *cobra.Command, client metadata.ProvisionWatcherClient, id string) ([]models.ProvisionWatcher, error) {
	watcher, err := client.ProvisionWatcher(cmd.Context(), id)
	if err != nil {
		return nil, err
	}
	return []models.ProvisionWatcher{watcher}, nil
}
