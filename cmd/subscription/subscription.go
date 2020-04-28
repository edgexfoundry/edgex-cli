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

package subscription

import (
	addsubscription "github.com/edgexfoundry-holding/edgex-cli/cmd/subscription/add"
	listsubscription "github.com/edgexfoundry-holding/edgex-cli/cmd/subscription/list"
	rmsubscription "github.com/edgexfoundry-holding/edgex-cli/cmd/subscription/rm"

	"github.com/spf13/cobra"
)

//deprecated
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "subscription",
		Short: "Subscription command",
		Long:  `Actions related to subscriptions.`,
	}
	cmd.AddCommand(addsubscription.NewCommand())
	cmd.AddCommand(rmsubscription.NewCommand())
	cmd.AddCommand(listsubscription.NewCommand())
	return cmd
}
