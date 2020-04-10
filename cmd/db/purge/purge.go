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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

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
		Run: func(cmd *cobra.Command, args []string) {

			// asking for user input
			fmt.Println("Are you sure? This cannot be undone: [y/n]")
			if askForConfirmation() {
				fmt.Println("Removing all objects from DB...")
				purge()
			} else {
				fmt.Println("Aborting")
				return
			}

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
	fmt.Println("* core-metadata")

	//////////////////////////////////////////////////////
	// DEVICE
	//////////////////////////////////////////////////////
	ctx, _ := context.WithCancel(context.Background())

	url := config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceRoute
	mdc := metadata.NewDeviceClient(local.New(url))

	devices, err := mdc.Devices(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	numberDevices := len(devices)
	for _, device := range devices {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + device.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Removed ", numberDevices, " devices.")

	//////////////////////////////////////////////////////
	// DS
	//////////////////////////////////////////////////////

	url = config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute
	deviceServiceData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}
	var deviceServices []models.DeviceService

	err = json.Unmarshal(deviceServiceData, &deviceServices)
	if err != nil {
		fmt.Println(err)
	}

	numberDSs := len(deviceServices)
	for _, deviceService := range deviceServices {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + deviceService.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Removed ", numberDSs, " device services.")

	//////////////////////////////////////////////////////
	// DP
	//////////////////////////////////////////////////////
	url = config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceProfileRoute
	deviceProfileData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var deviceProfiles []models.DeviceProfile

	err = json.Unmarshal(deviceProfileData, &deviceProfiles)
	if err != nil {
		fmt.Println(err)
	}

	numberProfiles := len(deviceProfiles)
	for _, deviceProfile := range deviceProfiles {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + deviceProfile.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Removed ", numberProfiles, " device profiles.")

	//////////////////////////////////////////////////////
	// ADDRESSABLES
	//////////////////////////////////////////////////////
	// Calling GetAllItems function, which
	// makes API call to get all items of given typ
	url = config.Conf.Clients["Metadata"].Url() + clients.ApiAddressableRoute
	data, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	// unmarshalling the json response
	var addressables []models.Addressable
	err = json.Unmarshal(data, &addressables)
	if err != nil {
		fmt.Println(err)
	}

	// Looping over the list of items and calling
	// DeleteItem for each

	numberItems := len(addressables)
	for _, addr := range addressables {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + addr.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Removed ", numberItems, " addressables.")

	// CORE-DATA:
	fmt.Println("* core-data")
	//////////////////////////////////////////////////////
	// Events and Readings
	//////////////////////////////////////////////////////

	removeEvents()

	//////////////////////////////////////////////////////
	// reading
	//////////////////////////////////////////////////////
	url = config.Conf.Clients["CoreData"].Url() + clients.ApiReadingRoute
	readingData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var readings []models.Reading

	err = json.Unmarshal(readingData, &readings)
	if err != nil {
		fmt.Println(err)
	}

	numberRs := len(readings)
	for _, reading := range readings {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + reading.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Removed ", numberRs, " readings.")

	//////////////////////////////////////////////////////
	// value descriptors
	//////////////////////////////////////////////////////
	url = config.Conf.Clients["CoreData"].Url() + clients.ApiValueDescriptorRoute
	valueDescriptorData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var valueDescriptors []models.ValueDescriptor

	err = json.Unmarshal(valueDescriptorData, &valueDescriptors)
	if err != nil {
		fmt.Println(err)
	}

	numberVDs := len(valueDescriptors)
	for _, valueDescriptor := range valueDescriptors {
		// call delete function here
		_, err = client.DeleteItem(url + config.PathId + valueDescriptor.Id)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Removed ", numberVDs, " value descriptors.")

	//////////////////////////////////////////////////////
	// Logs
	//////////////////////////////////////////////////////

	fmt.Println("* Logs")
	removeLogs()

	//////////////////////////////////////////////////////
	// Scheduler
	//////////////////////////////////////////////////////
	fmt.Println("* Scheduler")

	//////////////////////////////////////////////////////
	// Interval
	//////////////////////////////////////////////////////
	url = config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalRoute
	intervalData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var intervals []models.Interval

	err = json.Unmarshal(intervalData, &intervals)
	if err != nil {
		fmt.Println(err)
		return
	}

	numberIs := len(intervals)
	for _, interval := range intervals {

		// call delete function here
		_, err = client.DeleteItem(url + "/" + interval.ID)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Removed ", numberIs, " interval.")

	//////////////////////////////////////////////////////
	// Interval Action
	//////////////////////////////////////////////////////

	url = config.Conf.Clients["Scheduler"].Url() + clients.ApiIntervalActionRoute
	intervalActionData, err := client.GetAllItems(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var intervalActions []models.IntervalAction

	err = json.Unmarshal(intervalActionData, &intervalActions)
	if err != nil {
		fmt.Println(err)
	}

	numberIAs := len(intervalActions)
	for _, intervalAction := range intervalActions {
		// call delete function here
		_, err = client.DeleteItem(url + "/" + intervalAction.ID)

		if err != nil {
			fmt.Println(err)
			//TODO should we stop the execution of the purge command. Anywhy previous successful request cannot be reverted?
			return
		}
	}

	fmt.Println("Removed ", numberIAs, " interval action.")

	//////////////////////////////////////////////////////
	// notifications
	//////////////////////////////////////////////////////
	fmt.Println("* Notifications")
	removeNotifications()

	////////////////////////////////////////////////////////
	//// exportclient
	////////////////////////////////////////////////////////
	//
	//type registrationList struct {
	//	list []models.Registration
	//}
	//
	//registrationData, err := client.GetAllItemsDepricated("registration", "48071")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//registrations := registrationList{}
	//
	//registrationerrjson := json.Unmarshal(registrationData, &registrations.list)
	//if registrationerrjson != nil {
	//	fmt.Println(registrationerrjson)
	//}
	//
	//numberRegs := len(registrations.list)
	//for _, object := range registrations.list {
	//
	//	// call delete function here
	//	_, err = client.DeleteItemByIdOrName(object.ID,
	//		config.Conf.ExportService.RegistrationByIDRoute,
	//		config.Conf.ExportService.RegistrationByNameRoute,
	//		config.Conf.ExportService.Port)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}
	//
	//fmt.Println("Removed ", numberRegs, " registrations.")
}

func removeEvents() {

	// Create client
	client := &http.Client{}
	url := config.Conf.Clients["CoreData"].Url()+clients.ApiEventRoute+"/scruball"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, errResp := client.Do(req)
	if errResp != nil {
		fmt.Println(errResp)
		return
	}

	defer resp.Body.Close()

	respBody, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println(errBody)
		return
	}
	fmt.Printf("Removed %v events\n", string(respBody))
}

func removeLogs() {
	ts := time.Now().Unix() * 1000

	// Create client
	client := &http.Client{}
	url:=config.Conf.Clients["Logging"].Url()+clients.ApiLoggingRoute+"/0"
	req, err := http.NewRequest("DELETE", url+strconv.FormatInt(ts, 10), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, errResp := client.Do(req)
	if errResp != nil {
		fmt.Println(errResp)
		return
	}

	respBody, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println(errBody)
		return
	}

	defer resp.Body.Close()

	fmt.Printf("Removed %v logs\n", string(respBody))
}

func removeNotifications() {
	// Create client
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", config.Conf.Clients["Notification"].Url()+"/api/v1/cleanup", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, errResp := client.Do(req)
	if errResp != nil {
		fmt.Println(errResp)
		return
	}

	respBody, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println(errBody)
		return
	}

	defer resp.Body.Close()

	fmt.Printf("Removed notifications %v\n", string(respBody))
}
