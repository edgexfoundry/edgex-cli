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

package interval

import (
	addinterval "github.com/edgexfoundry-holding/edgex-cli/cmd/interval/add"
	listinterval "github.com/edgexfoundry-holding/edgex-cli/cmd/interval/list"
	rminterval "github.com/edgexfoundry-holding/edgex-cli/cmd/interval/rm"
	updateinterval "github.com/edgexfoundry-holding/edgex-cli/cmd/interval/update"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "interval",
		Short: "Interval command",
		Long:  `Actions related to intervals (scheduler).`,
	}
	cmd.AddCommand(addinterval.NewCommand())
	cmd.AddCommand(rminterval.NewCommand())
	cmd.AddCommand(updateinterval.NewCommand())
	cmd.AddCommand(listinterval.NewCommand())
	return cmd
}
