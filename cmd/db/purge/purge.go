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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// NewCommand returns the db command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purges entire EdgeX Database. [USE WITH CAUTION]",
		Long: `Purge DB

USE WITH CAUTION. The effect of this command is irreversible.

TODO: clean other tables:
-	coredata
-	logging
-	notifications
-	exportclient

The end goal for this command is to act like a clean_mongo.js script for any underlying 
database. Currently, it only cleans up core-metadata.


`,
		Run: func(cmd *cobra.Command, args []string) {

			verbose, _ := cmd.Flags().GetBool("verbose")

			fmt.Println("* core-metadata")

			//////////////////////////////////////////////////////
			// DEVICE
			//////////////////////////////////////////////////////
			type deviceList struct {
				list []models.Device
			}

			deviceData, err := client.GetAllItems("device", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			devices := deviceList{}
			deviceerrjson := json.Unmarshal(deviceData, &devices.list)
			if deviceerrjson != nil {
				fmt.Println(deviceerrjson)
			}

			numberDevices := len(devices.list)

			for _, object := range devices.list {
				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.MetadataService.DeviceServiceByIDRoute, config.Conf.MetadataService.DeviceServiceBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Println("Removed ", numberDevices, " devices.")

			//////////////////////////////////////////////////////
			// DR
			//////////////////////////////////////////////////////

			type deviceReportList struct {
				list []models.DeviceReport
			}

			deviceReportData, err := client.GetAllItems("devicereport", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			devicereports := deviceReportList{}

			devicereporterrjson := json.Unmarshal(deviceReportData, &devicereports.list)
			if devicereporterrjson != nil {
				fmt.Println(devicereporterrjson)
			}

			numberDRs := len(devicereports.list)
			for _, object := range devicereports.list {
				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.MetadataService.DeviceServiceByIDRoute, config.Conf.MetadataService.DeviceServiceBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Println("Removed ", numberDRs, " device reports.")

			//////////////////////////////////////////////////////
			// DS
			//////////////////////////////////////////////////////
			type deviceServiceList struct {
				list []models.DeviceService
			}

			deviceServiceData, err := client.GetAllItems("deviceservice", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}
			deviceservices := deviceServiceList{}

			deviceserviceerrjson := json.Unmarshal(deviceServiceData, &deviceservices.list)
			if deviceserviceerrjson != nil {
				fmt.Println(deviceserviceerrjson)
			}

			numberDSs := len(deviceservices.list)
			for _, object := range deviceservices.list {
				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.MetadataService.DeviceServiceByIDRoute, config.Conf.MetadataService.DeviceServiceBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberDSs, " device services.")

			//////////////////////////////////////////////////////
			// DP
			//////////////////////////////////////////////////////
			type deviceProfileList struct {
				list []models.DeviceProfile
			}

			DeviceProfileData, err := client.GetAllItems("deviceprofile", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			deviceprofiles := deviceProfileList{}

			deviceprofileerrjson := json.Unmarshal(DeviceProfileData, &deviceprofiles.list)
			if deviceprofileerrjson != nil {
				fmt.Println(deviceprofileerrjson)
			}

			numberProfiles := len(deviceprofiles.list)
			for _, object := range deviceprofiles.list {
				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.MetadataService.DeviceProfileByIDRoute, config.Conf.MetadataService.DeviceProfileBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Println("Removed ", numberProfiles, " device profiles.")

			//////////////////////////////////////////////////////
			// ADDRESSABLES
			//////////////////////////////////////////////////////

			type addressableList struct {
				list []models.Addressable
			}

			// Calling GetAllItems function, which
			// makes API call to get all items of given typ
			data, err := client.GetAllItems("addressable", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			// unmarshalling the json response
			list := addressableList{}
			errjson := json.Unmarshal(data, &list.list)
			if errjson != nil {
				fmt.Println(errjson)
			}

			// Looping over the list of items and calling
			// DeleteItem for each

			numberItems := len(list.list)
			for _, addr := range list.list {
				// call delete function here
				_, err = client.DeleteItem(addr.Id, config.Conf.MetadataService.DeviceServiceByIDRoute, config.Conf.MetadataService.DeviceServiceBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}
			fmt.Println("Removed ", numberItems, " device provision watchers.")

			//////////////////////////////////////////////////////
			// Provision watchers
			//////////////////////////////////////////////////////

			type provisionWatcherList struct {
				list []models.ProvisionWatcher
			}

			provisionWatcherData, err := client.GetAllItems("provisionwatcher", "48081", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			provisionwatchers := provisionWatcherList{}

			provisionwatchererrjson := json.Unmarshal(provisionWatcherData, &provisionwatchers.list)
			if provisionwatchererrjson != nil {
				fmt.Println(provisionwatchererrjson)
			}

			numberPRs := len(provisionwatchers.list)
			for _, object := range provisionwatchers.list {
				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.MetadataService.DeviceServiceByIDRoute, config.Conf.MetadataService.DeviceServiceBySlugNameRoute, config.Conf.MetadataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberPRs, " device provision watchers.")

			// CORE-DATA:
			fmt.Println("* core-data")
			//////////////////////////////////////////////////////
			// Events and Readings
			//////////////////////////////////////////////////////

			removeEventsAndReadings()

			//////////////////////////////////////////////////////
			// reading
			//////////////////////////////////////////////////////

			type readingList struct {
				list []models.Reading
			}

			readingData, err := client.GetAllItems("reading", "48080", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			readings := readingList{}

			readingerrjson := json.Unmarshal(readingData, &readings.list)
			if readingerrjson != nil {
				fmt.Println(readingerrjson)
			}

			numberRs := len(readings.list)
			for _, object := range readings.list {

				// call delete function here

				_, err = client.DeleteItem(object.Id, config.Conf.DataService.ReadingByIDRoute, "", config.Conf.DataService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberRs, " readings.")

			//////////////////////////////////////////////////////
			// value descriptors
			//////////////////////////////////////////////////////

			type valueDescriptorList struct {
				list []models.ValueDescriptor
			}

			valueDescriptorData, err := client.GetAllItems("valuedescriptor", "48080", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			valuedescriptors := valueDescriptorList{}

			valuedescriptorerrjson := json.Unmarshal(valueDescriptorData, &valuedescriptors.list)
			if valuedescriptorerrjson != nil {
				fmt.Println(valuedescriptorerrjson)
			}

			numberVDs := len(valuedescriptors.list)
			for _, object := range valuedescriptors.list {

				// call delete function here
				_, err = client.DeleteItem(object.Id, config.Conf.DataService.VDescriptorByIDRoute, config.Conf.DataService.VDescriptorByNameRoute, config.Conf.DataService.Port, verbose)

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
			type intervalList struct {
				list []models.Interval
			}

			intervalData, err := client.GetAllItems("interval", "48085", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			intervals := intervalList{}

			intervalerrjson := json.Unmarshal(intervalData, &intervals.list)
			if intervalerrjson != nil {
				fmt.Println(intervalerrjson)
			}

			numberIs := len(intervals.list)
			for _, object := range intervals.list {

				// call delete function here
				_, err = client.DeleteItem(object.ID, config.Conf.SchedulerService.IntervalByIDRoute,
					config.Conf.SchedulerService.IntervalByNameSlugRoute, config.Conf.SchedulerService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberIs, " interval.")

			//////////////////////////////////////////////////////
			// Interval Action
			//////////////////////////////////////////////////////
			type intervalactionList struct {
				list []models.IntervalAction
			}

			intervalactionData, err := client.GetAllItems("intervalaction", "48085", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			intervalactions := intervalactionList{}

			intervalactionerrjson := json.Unmarshal(intervalactionData, &intervalactions.list)
			if intervalactionerrjson != nil {
				fmt.Println(intervalactionerrjson)
			}

			numberIAs := len(intervalactions.list)
			for _, object := range intervalactions.list {

				// call delete function here
				_, err = client.DeleteItem(object.ID, config.Conf.SchedulerService.IntervalActionByIDRoute,
					config.Conf.SchedulerService.IntervalActionByNameSlugRoute, config.Conf.SchedulerService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberIAs, " interval action.")

			//////////////////////////////////////////////////////
			// notifications
			//////////////////////////////////////////////////////
			fmt.Println("* Notifications")
			removeNotifications()

			//////////////////////////////////////////////////////
			// exportclient
			//////////////////////////////////////////////////////

			type registrationList struct {
				list []models.Registration
			}

			registrationData, err := client.GetAllItems("registration", "48071", verbose)

			if err != nil {
				fmt.Println(err)
				return
			}

			registrations := registrationList{}

			registrationerrjson := json.Unmarshal(registrationData, &registrations.list)
			if registrationerrjson != nil {
				fmt.Println(registrationerrjson)
			}

			numberRegs := len(registrations.list)
			for _, object := range registrations.list {

				// call delete function here
				_, err = client.DeleteItem(object.ID, config.Conf.ExportService.RegistrationByIDRoute,
					config.Conf.ExportService.RegistrationByNameRoute, config.Conf.ExportService.Port, verbose)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Removed ", numberRegs, " registrations.")

		},
	}
	return cmd
}

func removeEventsAndReadings() {
	host := viper.GetString("Host")

	// Create client
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://"+host+":48080/api/v1/event/scrub", nil)
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
	host := viper.GetString("Host")

	ts := time.Now().Unix() * 1000

	// Create client
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://"+host+":48061/api/v1/logs/0/"+strconv.FormatInt(ts, 10), nil)
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
	host := viper.GetString("Host")

	// Create client
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://"+host+":48060/api/v1/cleanup", nil)
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
