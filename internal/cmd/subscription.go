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
	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:          "subscription",
		Short:        "Add, remove and list subscriptions [Support Notificationss]",
		Long:         "Add, remove and list subscriptions [Support Notifications]",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	initRmSubscriptionCommand(cmd)
	initAddSubscriptionCommand(cmd)
	initListSubscriptionCommand(cmd)
	initGetSubscriptionByNameCommand(cmd)
}

// initRmSubscriptionCommand implements the DELETE /subscription/name/{name} endpoint
// "Deletes a subscription according to the given name."
func initRmSubscriptionCommand(cmd *cobra.Command) {
	var rm = &cobra.Command{
		Use:          "rm",
		Short:        "Delete the named subscription",
		Long:         "Delete the named subscription",
		RunE:         handleRmSubscription,
		SilenceUsage: true,
	}
	rm.Flags().StringVarP(&subscriptionName, "name", "n", "", "Name of subscription to remove")
	rm.MarkFlagRequired("name")
	cmd.AddCommand(rm)
}

var subscriptionName, subscriptionChannels, subscriptionReceiver string
var subscriptionCategories, subscriptionDescription, subscriptionResendInterval string
var subscriptionResendLimit int
var subscriptionAdminState, subscriptionSelectedCategory, subscriptionSelectedLabel, subscriptionSelectedReceiver string

func getSubscriptionChannels() (channels []dtos.Address, err error) {
	if subscriptionChannels != "" {
		err = jsonpkg.Unmarshal([]byte(subscriptionChannels), &channels)
		if err != nil {
			err = fmt.Errorf("channels JSON object array invalid (%v)", err)
		}
	}
	return
}

func getSubscriptionCategories() []string {
	var a []string
	if len(subscriptionCategories) > 0 {
		a = strings.Split(subscriptionCategories, ",")
	}
	return a

}

// initAddSubscriptionCommand implements the POST /notification endpoint
// "Adds one or more notifications to be sent."
func initAddSubscriptionCommand(cmd *cobra.Command) {
	var add = &cobra.Command{
		Use:          "add",
		Short:        "Add a new subscription",
		Long:         "Add a new subscription",
		Example:      `  edgex-cli subscription add -n "name01" --receiver "receiver01" -c "[{\"type\": \"REST\", \"host\": \"localhost\", \"port\": 7770, \"httpMethod\": \"POST\"}]"`,
		RunE:         handleAddSubscription,
		SilenceUsage: true,
	}

	add.Flags().StringVarP(&subscriptionName, "name", "n", "", "A meaningful identifier for the subscription")
	add.Flags().StringVarP(&subscriptionChannels, "channels", "c", "", "A JSON object array indicating how this subscription is capable of receiving notifications")
	add.Flags().StringVarP(&subscriptionCategories, "categories", "", "", "A comma-delimited list of categories")
	add.Flags().StringVarP(&subscriptionReceiver, "receiver", "", "", "The name of the party interested in the notification")
	add.Flags().StringVarP(&subscriptionDescription, "description", "", "", "An optional description of the subscription's intent.")
	add.Flags().IntVarP(&subscriptionResendLimit, "resend-limit", "", 0, "The retry limit for attempts to send notifications")
	add.Flags().StringVarP(&subscriptionResendInterval, "resend-interval", "", "1h", "The interval in ISO 8691 format of resending the notification")
	add.Flags().StringVarP(&subscriptionAdminState, "admin-state", "a", "UNLOCKED", "Admin state [LOCKED | UNLOCKED]")

	addLabelsFlag(add)
	add.MarkFlagRequired("name")
	add.MarkFlagRequired("receiver")
	add.MarkFlagRequired("channels")
	//admin-state"
	//channels
	cmd.AddCommand(add)
}

// initListSubscriptionCommand implements a number of endpoints:
// GET /subscription/all
// "Allows paginated retrieval of subscriptions, sorted by created timestamp descending."
// GET /subscription/category/{category}
// "Returns a paginated list of subscriptions associated with the specified category."
// GET /subscription/label/{label}
// "Returns a paginated list of subscriptions associated with the specified label."
// GET /subscription/receiver/{receiver}
// "Returns a paginated list of subscriptions associated with the specified receiver."
func initListSubscriptionCommand(cmd *cobra.Command) {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List subscriptions",
		Long:  "List all subscriptions, optionally filtered by a given category, label or receiver",
		Example: `  edgex-cli subscription list
  edgex-cli subscription list --category "my-category"
  edgex-cli subscription list --label "my-label"`,
		RunE:         handleListSubscription,
		SilenceUsage: true,
	}
	listCmd.Flags().StringVarP(&subscriptionSelectedCategory, "category", "c", "", "List subscriptions associated with this category")
	listCmd.Flags().StringVarP(&subscriptionSelectedLabel, "label", "", "", "List subscriptions associated with this label")
	listCmd.Flags().StringVarP(&subscriptionSelectedReceiver, "receiver", "r", "", "List subscriptions associated with this receiver")
	addFormatFlags(listCmd)
	addVerboseFlag(listCmd)
	addLimitOffsetFlags(listCmd)
	cmd.AddCommand(listCmd)
}

// initGetSubscriptionByNameCommand implements the GET ​/subscription/name endpoint
// "Returns a subscription by its unique name.""
func initGetSubscriptionByNameCommand(cmd *cobra.Command) {
	var nameCmd = &cobra.Command{
		Use:          "name",
		Short:        "Return a subscription by its unique name",
		Long:         `Return a subscription by its unique name`,
		RunE:         handleGetSubscriptionByName,
		SilenceUsage: true,
	}
	nameCmd.Flags().StringVarP(&subscriptionName, "name", "n", "", "Subscription name")
	nameCmd.MarkFlagRequired("name")
	addFormatFlags(nameCmd)
	addVerboseFlag(nameCmd)
	cmd.AddCommand(nameCmd)

}

func handleRmSubscription(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetSubscriptionClient()
	response, err := client.DeleteSubscriptionByName(context.Background(), subscriptionName)
	if err == nil {
		fmt.Println(response.Message)
	}
	return err
}

func handleAddSubscription(cmd *cobra.Command, args []string) error {

	client := getSupportNotificationsService().GetSubscriptionClient()

	err := validateAdminState(subscriptionAdminState)
	if err != nil {
		return err
	}

	channels, err := getSubscriptionChannels()
	if err != nil {
		return err
	}

	l := getLabels()
	c := getSubscriptionCategories()

	if l == nil && c == nil {
		return errors.New("either labels or categories must be specified")

	}

	var req = requests.NewAddSubscriptionRequest(dtos.Subscription{
		Name:           subscriptionName,
		Channels:       channels,
		Receiver:       subscriptionReceiver,
		Categories:     c,
		Labels:         l,
		Description:    subscriptionDescription,
		ResendLimit:    subscriptionResendLimit,
		ResendInterval: subscriptionResendInterval,
		AdminState:     subscriptionAdminState,
	})

	response, err := client.Add(context.Background(), []requests.AddSubscriptionRequest{req})

	if err != nil {
		return err
	}
	if response != nil {
		fmt.Println(response[0])
	}
	return err
}

func handleListSubscription(cmd *cobra.Command, args []string) error {

	client := getSupportNotificationsService().GetSubscriptionClient()

	var response responses.MultiSubscriptionsResponse
	var err error

	if subscriptionSelectedCategory != "" {
		response, err = client.SubscriptionsByCategory(context.Background(), subscriptionSelectedCategory, offset, limit)
	} else if subscriptionSelectedLabel != "" {
		response, err = client.SubscriptionsByLabel(context.Background(), subscriptionSelectedLabel, offset, limit)
	} else if subscriptionSelectedReceiver != "" {
		response, err = client.SubscriptionsByReceiver(context.Background(), subscriptionSelectedReceiver, offset, limit)
	} else {
		response, err = client.AllSubscriptions(context.Background(), offset, limit)
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

		if len(response.Subscriptions) == 0 {
			fmt.Println("No subscriptions available")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		printSubscriptionTableHeader(w)
		for _, n := range response.Subscriptions {
			printSubscription(w, &n)
		}
		w.Flush()
	}
	return nil
}

func handleGetSubscriptionByName(cmd *cobra.Command, args []string) error {
	client := getSupportNotificationsService().GetSubscriptionClient()

	response, err := client.SubscriptionByName(context.Background(), subscriptionName)
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
		printSubscriptionTableHeader(w)
		printSubscription(w, &response.Subscription)
		w.Flush()
	}
	return nil
}

func printSubscriptionTableHeader(w *tabwriter.Writer) {
	if verbose {
		fmt.Fprintln(w, "Id\tName\tDescription\tChannels\tReceiver\tCategories\tLabels\tResendLimit\tResendInterval\tAdminState")
	} else {
		fmt.Fprintln(w, "Namet\tDescription\tChannels\tReceiver\tCategories\tLabels")
	}
}

func printSubscription(w *tabwriter.Writer, n *dtos.Subscription) {
	if verbose {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			n.Id,
			n.Name,
			n.Description,
			n.Channels,
			n.Receiver,
			n.Categories,
			n.Labels,
			n.ResendLimit,
			n.ResendInterval,
			n.AdminState)
	} else {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
			n.Name,
			n.Description,
			n.Channels,
			n.Receiver,
			n.Categories,
			n.Labels)
	}
}
