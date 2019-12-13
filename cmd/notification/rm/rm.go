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
var byAge bool

func removeNotificationHandler(cmd *cobra.Command, args []string) {

	// Checking for args
	if len(args) == 0 {
		fmt.Printf("Error: No profile ID/Name provided.\n")
		return
	}

	var url string
	if byAge {
		url = config.Conf.NotificationService.NotificationByAgeRoute
	} else {
		url = config.Conf.NotificationService.NotificationByNameSlugRoute
	}

	respBody, err := client.DeleteItemByName(args[0],
		url,
		config.Conf.MetadataService.Port)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Display Results
	if string(respBody) == "true" {
		fmt.Printf("Removed: %s\n", args[0])
	} else {
		fmt.Printf("Remove Unsuccessful: %s\n", respBody)
	}
}

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [notification slug or age]",
		Short: "Removes notification by slug or age",
		Long:  `Removes a notification given its slug or age timestamp.`,
		Args:  cobra.ExactValidArgs(1),
		Run:   removeNotificationHandler,
	}
	cmd.Flags().BoolVar(&byAge, "age", false, "Remove by age")
	return cmd
}
