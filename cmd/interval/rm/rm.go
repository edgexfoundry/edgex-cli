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
	"context"
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/scheduler"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/cobra"
)

func removeIntervalHandler(cmd *cobra.Command, args []string) {
	// Checking for args
	if len(args) == 0 {
		fmt.Printf("Error: No interval ID/Name provided.\n")
		return
	}

	url := config.Conf.Clients["Scheduler"].Url()
	// Create request
	intervalID := args[0]
	sc := scheduler.NewIntervalClient(
		local.New(url+clients.ApiIntervalRoute),
	)

	err := sc.Delete(context.Background(), intervalID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Removed: %s\n", intervalID)
}

// NewCommand returns rm interval command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [interval name or id]",
		Short: "Removes interval by name or id",
		Long:  `Removes an interval given its name or id.`,
		Args:  cobra.ExactValidArgs(1),
		Run:   removeIntervalHandler,
	}
	return cmd
}
