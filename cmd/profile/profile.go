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

package profile

import (
	listprofile "github.com/edgexfoundry/edgex-cli/cmd/profile/list"
	rmprofile "github.com/edgexfoundry/edgex-cli/cmd/profile/rm"
	addprofile "github.com/edgexfoundry/edgex-cli/cmd/profile/add"
	"github.com/spf13/cobra"
)

// NewCommand returns the profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile [option]",
		Short: "Device profile command.",
		Long:  `Actions associated with device profiles`,
	}

	cmd.AddCommand(rmprofile.NewCommand())
	cmd.AddCommand(listprofile.NewCommand())
	cmd.AddCommand(addprofile.NewCommand())
	return cmd
}
