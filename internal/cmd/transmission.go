/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package cmd

import (
	"context"
	jsonpkg "encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/responses"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/spf13/cobra"
)

var transmissionId string
var transmissionAge int
var transmissionSubscriptionName, transmissionStart, transmissionEnd, transmissionStatus string

func init() {
	var cmd = &cobra.Command{
		Use:          "transmission",
		Short:        "Remove and list transmissions [Support Notifications]",
		Long:         "Remove and list transmissions [Support Notifications]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initRmTransmissionCommand(cmd)
	initListTransmissionCommand(cmd)
	initGetTransmissionByIdCommand(cmd)
}

// initRmTransmissionCommand implements the DELETE /transmission/age/{age}
// "Deletes the processed transmissions if the current timestamp minus
// their created timestamp is less than the age parameter."
func initRmTransmissionCommand(cmd *cobra.Command) {
	var rm = &cobra.Command{
		Use:          "rm",
		Short:        "Delete processed transmissions",
		Long:         "Delete processed transmissions older than the specificed age (in milliseconds)",
		RunE:         handleRmTransmission,
		SilenceUsage: true,
	}
	rm.Flags().IntVarP(&transmissionAge, "age", "a", 0, "The minimum age of transmissions to deleted (in milliseconds)")
	rm.MarkFlagRequired("age")
	cmd.AddCommand(rm)
}

// initListTransmissionCommand implements a number of endpoints:
// - GET /transmission/all
//   "Given the entire range of transmissions sorted in descending order of created time,
//   returns a portion of that range according to the offset and limit parameters."
// - GET /transmission/subscription/name/{name}
//   "Returns a paginated list of transmissions that originated with the specified subscription."
// - GET /transmission/start/{start}/end/{end}
//   "Allows querying of transmissions by their creation timestamp within a
//   given time range, sorted in descending order. Results are paginated.
// - GET /transmission/status/{status}
//   "Allows retrieval of the transmissions associated with the specified status.
//    Ordered by create timestamp descending.""
func initListTransmissionCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List transmissions",
		Long:  "Get paginated list of transmissions, optionally filtered by a subscription name, status or time range",
		Example: `  edgex-cli transmission list
  edgex-cli transmission list --name "name01"
  edgex-cli transmission list --status "SENT"
  edgex-cli transmission list --start "01 jan 20 00:00 GMT" --end "01 dec 21 00:00 GMT"`,
		RunE:         handleListTransmission,
		SilenceUsage: true,
	}
	listCmd.Flags().StringVarP(&transmissionSubscriptionName, "name", "n", "", "List transmissions that originated with the specified subscription")
	listCmd.Flags().StringVarP(&transmissionStart, "start", "s", "", "List transmissions from after this (RFC822) timestamp")
	listCmd.Flags().StringVarP(&transmissionEnd, "end", "e", "", "List transmissions from before this (RFC822) timestamp")
	listCmd.Flags().StringVarP(&transmissionStatus, "status", "", "", "List transmissions with this status [ACKNOWLEDGED, FAILED, SENT, RESENDING, ESCALATED]")
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
}

// initGetTransmissionByIdCommand implements the GET ​/transmission/id endpoint
// "Returns a transmission by ID."
func initGetTransmissionByIdCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "id",
		Short:        "Return a transmission by ID",
		Long:         `Return a transmission by ID`,
		RunE:         handleGetTransmissionById,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&transmissionId, "id", "i", "", "The ID that identifies the transmission")
	nameCmd.MarkFlagRequired("id")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleRmTransmission(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetTransmissionClient()
	response, err := client.DeleteProcessedTransmissionsByAge(context.Background(), transmissionAge)
	if err == nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleListTransmission(cmd *cobra.Command, args []string) error {

	var transmissionSubscriptionName, transmissionStart, transmissionEnd, transmissionStatus string

	client := getSupportNotificationsService().GetTransmissionClient()
	var response responses.MultiTransmissionsResponse
	var err error

	if transmissionSubscriptionName != "" {
		response, err = client.TransmissionsBySubscriptionName(context.Background(), transmissionSubscriptionName, offset, limit)
	} else if transmissionStatus != "" {
		transmissionStatus = strings.ToUpper(transmissionStatus)
		if !(transmissionStatus == models.Acknowledged || notificationStatus == models.Failed || notificationStatus == models.Sent ||
			notificationStatus == models.RESENDING || notificationStatus == models.Escalated) {
			return fmt.Errorf("status should be one of: %s, %s, %s, %s, %s", models.Acknowledged, models.Failed, models.Sent,
				models.RESENDING, models.Escalated)
		}
		response, err = client.TransmissionsByStatus(context.Background(), transmissionStatus, offset, limit)
	} else if transmissionStart != "" && transmissionEnd != "" {
		start, err := getMillisTimestampFromRFC822Time(transmissionStart)
		if err != nil {
			return err
		}
		end, err := getMillisTimestampFromRFC822Time(transmissionEnd)
		if err != nil {
			return err
		}
		response, err = client.TransmissionsByTimeRange(context.Background(), int(start), int(end), offset, limit)
	} else {
		response, err = client.AllTransmissions(context.Background(), offset, limit)

	}

	if err != nil {
		return err
	}

	if json {
		result, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}
		fmt.Print(string(result))
	} else {

		if len(response.Transmissions) == 0 {
			fmt.Println("No transmissions available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printTransmissionTableHeader(w)
		for _, n := range response.Transmissions {
			printTransmission(w, &n)
		}
		w.Flush()
	}
	return nil
}

func handleGetTransmissionById(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetTransmissionClient()

	response, err := client.TransmissionById(context.Background(), transmissionId)
	if err != nil {
		return err
	}

	if json {
		result, err := jsonpkg.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Println(string(result))
	} else {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printTransmissionTableHeader(w)
		printTransmission(w, &response.Transmission)
		w.Flush()
	}
	return nil
}

func printTransmissionTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tChannel\tCreated\tNotificationId\tSubscriptionName\tRecords\tResendCount\tStatus")
	} else {
		fmt.Fprintln(w, "SubscriptionName\tResendCount\tStatus")
	}

}

func printTransmission(w *tabwriter.Writer, t *dtos.Transmission) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			t.Id,
			t.Channel,
			getRFC822Time(t.Created),
			t.NotificationId,
			t.SubscriptionName,
			t.Records,
			t.ResendCount,
			t.Status)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\n",
			t.SubscriptionName,
			t.ResendCount,
			t.Status)

	}
}
