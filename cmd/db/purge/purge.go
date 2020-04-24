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
	"log"
	"strconv"
	"time"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
	cleaners "github.com/edgexfoundry-holding/edgex-cli/pkg/cmd/purge"

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

The end goal for this command is to act like a clean_mongo.js script for any underlying 
database. Currently, it only cleans up core-metadata.


`,
		RunE: func(cmd *cobra.Command, args []string) (err error){
			// asking for user input
			fmt.Println("Are you sure? This cannot be undone: [y/n]")
			if askForConfirmation() {
				fmt.Println("Removing all objects from DB...")
				purge()
			} else {
				fmt.Println("Aborting")
				return
			}

			return
		},
	}
	return cmd
}

// three following functions where found here: https://gist.github.com/albrow/5882501

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Are you sure? This cannot be undone:")
		return askForConfirmation()
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
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
	_, err := client.DeleteItem(url)
	if err != nil {
		fmt.Println(err)
	} else {
		//TODO fix the message
		fmt.Print("Logs removed\n")
	}
}

func removeNotifications(){
	fmt.Println("\n * Notifications")
	url := config.Conf.Clients["Notification"].Url()+"/api/v1/cleanup"
	_, err := client.DeleteItem(url)
	if err != nil {
		fmt.Println(err)
	} else {
		//TODO fix the message
		fmt.Println("All Notification have been removed")
	}
}
