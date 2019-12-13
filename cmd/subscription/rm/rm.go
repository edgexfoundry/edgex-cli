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

// NewCommand returns the rm command of type cobra.Command
var byID bool

func removeSubscriptionHandler(cmd *cobra.Command, args []string) {
	// Checking for args
	if len(args) == 0 {
		fmt.Printf("Error: No subription ID/Name provided.\n")
		return
	}

	// Create request
	subscriptionlID := args[0]

	respBody, err := client.DeleteItem(subscriptionlID,
		config.Conf.NotificationService.SubscriptionByIDRoute,
		config.Conf.NotificationService.SubscriptionByNameSlugRoute,
		config.Conf.NotificationService.Port)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Display Results
	if string(respBody) == "true" {
		fmt.Printf("Removed: %s\n", subscriptionlID)
	} else {
		fmt.Printf("Remove Unsuccessful!\n")
	}
}

// NewCommand returns remove subscription command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [notification slug or age]",
		Short: "Removes notification by slug or age",
		Long:  `Removes a notification given its slug or age timestamp.`,
		Args:  cobra.ExactValidArgs(1),
		Run:   removeSubscriptionHandler,
	}
	cmd.Flags().BoolVar(&byID, "id", false, "Remove by id")
	return cmd
}
