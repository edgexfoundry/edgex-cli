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

// Package purgedb command purges the entire Database. It performs the same action as the
// clean_mongo.js developer script. Unlike the clean_mongo.js, this command purges the
// database using API calls only. clean_mongo.js accesses the DB directly, which might
// always be possible using the CLI.
package purgedb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"
	cleaners "github.com/edgexfoundry/edgex-cli/pkg/cmd/purge"
	"github.com/edgexfoundry/edgex-cli/pkg/confirmation"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

// NewCommand returns the db command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purges entire EdgeX Database. [USE WITH CAUTION]",
		Long: `Purge DB

USE WITH CAUTION. The effect of this command is irreversible.

The end goal for this command is to clean all data from the underlying 
database.

`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// asking user to confirm the purge action
			if !confirmation.New().Confirm() {
				return
			}
			purge()
			return
		},
	}
	return cmd
}

func purge() {
	cleaners.NewMetadataCleaner().Purge()
	cleaners.NewCoredataCleaner().Purge()
	removeLogs()
	cleaners.NewSchedulerCleaner().Purge()
	removeNotifications()
}

func removeLogs() {
	fmt.Println("\n * Logs")
	ts := time.Now().Unix() * 1000
	url := config.Conf.Clients["Logging"].Url() + clients.ApiLoggingRoute + "/0/" + strconv.FormatInt(ts, 10)
	err := request.Delete(url)
	if err == nil {
		//TODO fix the message
		fmt.Print("Logs have been removed\n")
	}
}

func removeNotifications() {
	fmt.Println("\n * Notifications")
	url := config.Conf.Clients["Notification"].Url() + "/api/v1/cleanup"
	err := request.Delete(url)
	if err == nil {
		//TODO fix the message
		fmt.Println("Notifications have been removed")
	}
}
