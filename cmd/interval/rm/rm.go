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
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/spf13/cobra"
)

func removeIntervalHandler(cmd *cobra.Command, args []string) {
	// Create request
	verbose, _ := cmd.Flags().GetBool("verbose")
	intervalID := args[0]
	respBody := client.DeleteItem(intervalID, config.Conf.SchedulerService.IntervalByIDRoute,
		config.Conf.SchedulerService.IntervalByNameSlugRoute, config.Conf.SchedulerService.Port, verbose)

	// Display Results
	if string(respBody) == "true" {
		fmt.Printf("Removed: %s\n", intervalID)
	} else {
		fmt.Printf("Remove Unsuccessful!\n")
	}
}

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
