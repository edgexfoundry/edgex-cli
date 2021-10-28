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
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/responses"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:          "notification",
		Short:        "Add, remove and list notifications [Support Notifications]",
		Long:         "Add, remove and list notifications [Support Notifications]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initAddNotificationCommand(cmd)
	initListNotificationCommand(cmd)
	initRmNotificationCommand(cmd)
	initCleanupNotificationCommand(cmd)

}

var notificationCategory, notificationContent, notificationContentType, notificationDescription string
var notificationSender, notificationSeverity, notificationStatus string
var notificationLabel, notificationStart, notificationEnd, notificationId string

// initCleanupNotificationCommand implements DELETE /cleanup
// "Deletes all notifications and the corresponding transmissions.""
func initCleanupNotificationCommand(cmd *cobra.Command) {
	var cleanup = &cobra.Command{
		Use:          "cleanup",
		Short:        "Delete all notifications and corresponding transmissions",
		Long:         "Delete all notifications and corresponding transmissions",
		RunE:         handleCleanupNotifications,
		SilenceUsage: true,
	}
	cmd.AddCommand(cleanup)
}

// initRmDeviceCommand implements the DELETE /notification/id/{id}
// "Deletes a notification by ID and all of its associated transmissions.""
func initRmNotificationCommand(cmd *cobra.Command) {
	var rm = &cobra.Command{
		Use:          "rm",
		Short:        "Delete a notification and all of its associated transmissions",
		Long:         "Delete a notification and all of its associated transmissions",
		RunE:         handleRmNotifications,
		SilenceUsage: true,
	}
	rm.Flags().StringVarP(&notificationId, "id", "i", "", "The ID that identifies the notification")
	rm.MarkFlagRequired("id")
	cmd.AddCommand(rm)
}

// initAddNotificationCommand implements the POST /notification endpoint
// "Adds one or more notifications to be sent."
func initAddNotificationCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "add",
		Short:        "Add a notification to be sent",
		Long:         "Add a notification to be sent",
		Example:      `  edgex-cli notification add -s "sender01" -c "content" --category "category01"`,
		RunE:         handleAddNotifications,
		SilenceUsage: true,
	}
	add.Flags().StringVarP(&notificationCategory, "category", "", "", "Categorizes the notification")
	add.Flags().StringVarP(&notificationContent, "content", "c", "", "The content to be sent as the body of the notification")
	add.Flags().StringVarP(&notificationContentType, "content-type", "t", "", "Indicates the MIME type/Content-type of the notification's content")
	add.Flags().StringVarP(&notificationDescription, "description", "d", "", "An optional description of the notification's intent")
	add.Flags().StringVarP(&notificationSender, "sender", "s", "", "Identifies the sender of a notification, usually the name of sender")
	add.Flags().StringVarP(&notificationSeverity, "severity", "", "NORMAL", "Indicates the level of severity for the notification. Current accepted values include: MINOR, NORMAL, CRITICAL")
	add.Flags().StringVarP(&notificationStatus, "status", "", "", "A status indicating the current processing status of the notification. Accepted values are: NEW, PROCESSED, ESCALATED")
	addLabelsFlag(add)
	add.MarkFlagRequired("sender")
	add.MarkFlagRequired("content")
	add.MarkFlagRequired("category")
	cmd.AddCommand(add)
}

// initListNotificationCommand implements a number of endpoints:
// GET /notification/category/{category}
// "Returns a paginated list of notifications associated with the given category."
// GET /notification/label/{label}
// "Returns a paginated list of notifications associated with the given label."
// GET /notification/start/{start}/end/{end}
// "Allows querying of notifications by their creation timestamp within a given time range, sorted in descending order. Results are paginated.""
// GET /notification/status/{status}
// "Returns a paginated list of notifications with the specified status."
func initListNotificationCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List notifications",
		Long:  "List notifications associated with a given label, category or time range",
		Example: `  edgex-cli notification list --start "01 jan 20 00:00 GMT" --end "01 dec 21 00:00 GMT"
  edgex-cli notification list --category "category01"
  edgex-cli notification list --label "l01"`,
		RunE:         handleListNotifications,
		SilenceUsage: true,
	}
	listCmd.Flags().StringVarP(&notificationCategory, "category", "c", "", "List notifications belonging to this category")
	listCmd.Flags().StringVarP(&notificationLabel, "label", "", "", "List notifications with this label")
	listCmd.Flags().StringVarP(&notificationStart, "start", "s", "", "List notifications from after this (RFC822) timestamp")
	listCmd.Flags().StringVarP(&notificationEnd, "end", "e", "", "List notifications from before this (RFC822) timestamp")
	listCmd.Flags().StringVarP(&notificationStatus, "status", "", "", "List notifications with this status")

	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
}

func handleCleanupNotifications(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetNotificationClient()

	response, err := client.CleanupNotifications(context.Background())
	if err != nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleRmNotifications(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetNotificationClient()
	response, err := client.DeleteNotificationById(context.Background(), notificationId)
	if err == nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleAddNotifications(cmd *cobra.Command, args []string) error {

	if notificationStatus != "" {
		notificationStatus = strings.ToUpper(notificationStatus)
		if !(notificationStatus == models.New || notificationStatus == models.Processed || notificationStatus == models.Escalated) {
			return fmt.Errorf("status should be %s, %s or %s", models.New, models.Processed, models.Escalated)
		}
	}

	if notificationSeverity != "" {
		notificationSeverity = strings.ToUpper(notificationSeverity)
		if !(notificationSeverity == models.Minor || notificationSeverity == models.Normal || notificationSeverity == models.Critical) {
			return fmt.Errorf("severity should be %s, %s or %s", models.Minor, models.Normal, models.Critical)
		}
	}

	client := getSupportNotificationsService().GetNotificationClient()

	var req = requests.NewAddNotificationRequest(dtos.Notification{
		Category:    notificationCategory,
		Content:     notificationContent,
		ContentType: notificationContentType,
		Description: notificationDescription,
		Sender:      notificationSender,
		Severity:    notificationSeverity,
		Status:      notificationStatus,
		Labels:      getLabels(),
	})

	response, err := client.SendNotification(context.Background(), []requests.AddNotificationRequest{req})

	if err != nil {
		return err
	}
	if response != nil {
		fmt.Println(response[0])
	}
	return err
}

func handleListNotifications(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetNotificationClient()
	var response responses.MultiNotificationsResponse
	var err error

	if notificationCategory != "" {
		response, err = client.NotificationsByCategory(context.Background(), notificationCategory, offset, limit)
	} else if notificationLabel != "" {
		response, err = client.NotificationsByLabel(context.Background(), notificationLabel, offset, limit)
	} else if notificationStatus != "" {
		notificationStatus = strings.ToUpper(notificationStatus)
		if !(notificationStatus == models.New || notificationStatus == models.Processed || notificationStatus == models.Escalated) {
			return fmt.Errorf("status should be %s, %s or %s", models.New, models.Processed, models.Escalated)
		}
		response, err = client.NotificationsByStatus(context.Background(), notificationStatus, offset, limit)
	} else if notificationStart != "" && notificationEnd != "" {
		start, err := getMillisTimestampFromRFC822Time(notificationStart)
		if err != nil {
			return err
		}
		end, err := getMillisTimestampFromRFC822Time(notificationEnd)
		if err != nil {
			return err
		}
		response, err = client.NotificationsByTimeRange(context.Background(), int(start), int(end), offset, limit)
	} else {
		return errors.New("category, label, status or a timerange must be specified")
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

		if len(response.Notifications) == 0 {
			fmt.Println("No notifications available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printNotificationTableHeader(w)
		for _, n := range response.Notifications {
			printNotification(w, &n)
		}
		w.Flush()
	}
	return nil
}

func printNotificationTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tCategory\tContent\tContentType\tCreated\tDescription\tLabels\tModified\tSender\tSeverity\tStatus")
	} else {
		fmt.Fprintln(w, "Category\tContent\tDescription\tLabels\tSender\tSeverity\tStatus")
	}

}

func printNotification(w *tabwriter.Writer, n *dtos.Notification) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			n.Id,
			n.Category,
			n.Content,
			n.ContentType,
			getRFC822Time(n.Created),
			n.Description,
			n.Labels,
			getRFC822Time(n.Modified),
			n.Sender,
			n.Severity,
			n.Status)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			n.Category,
			n.Content,
			n.Description,
			n.Labels,
			n.Sender,
			n.Severity,
			n.Status)
	}
}
