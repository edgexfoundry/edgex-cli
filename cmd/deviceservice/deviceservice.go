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

package deviceservice

import (
	listdeviceservice "github.com/edgexfoundry/edgex-cli/cmd/deviceservice/list"
	rmdeviceservice "github.com/edgexfoundry/edgex-cli/cmd/deviceservice/rm"
	"github.com/spf13/cobra"
)

// NewCommand returns the device command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "deviceservice",
		Short: "Event command",
		Long:  `Actions related to device-generated events.`,
	}
	cmd.AddCommand(rmdeviceservice.NewCommand())
	cmd.AddCommand(listdeviceservice.NewCommand())
	return cmd
}
