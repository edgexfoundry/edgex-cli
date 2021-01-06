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

package rm

import (
	"errors"
	"fmt"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/scheduler"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

var name string

func removeIntervalHandler(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 && name == "" {
		return errors.New("no device id/name provided.\n")
	}

	sc := scheduler.NewIntervalClient(
		local.New(config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute),
	)
	var deletedBy string
	if name != "" {
		deletedBy = name
		err = sc.DeleteByName(cmd.Context(), name)
		if err != nil {
			return
		}
		fmt.Printf("Removed: %s\n", deletedBy)
	} else if len(args[0]) != 0 {
		return request.DeleteByIds(cmd.Context(), sc, args)
	}
	return nil
}

// NewCommand returns rm interval command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [ name | id]",
		Short: "Removes interval by name or id",
		Long:  `Removes an interval given its name or id.`,
		//Args:  cobra.ExactValidArgs(1),
		RunE: removeIntervalHandler,
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Delete interval by name")
	return cmd
}
