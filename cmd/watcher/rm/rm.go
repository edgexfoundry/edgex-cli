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
package rm

import (
	"errors"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

// NewCommand return the rm watcher command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [<id> ... ]",
		Short: "Remove watcher(s) by ID(s)",
		Long:  `Remove one or more watchers by given ID(s)`,
		RunE:  removeHandler,
	}
	return cmd
}

func removeHandler(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please provide watcher id(s).\n")
	}

	client := metadata.NewProvisionWatcherClient(
		local.New(config.Conf.Clients["Metadata"].Url() + clients.ApiProvisionWatcherRoute),
	)
	return request.DeleteByIds(cmd.Context(), client, args)
}
