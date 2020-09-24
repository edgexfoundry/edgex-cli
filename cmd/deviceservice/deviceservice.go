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
	"github.com/edgexfoundry/edgex-cli/cmd/deviceservice/add"
	"github.com/edgexfoundry/edgex-cli/cmd/deviceservice/list"
	"github.com/edgexfoundry/edgex-cli/cmd/deviceservice/rm"
	"github.com/edgexfoundry/edgex-cli/cmd/deviceservice/update"

	"github.com/spf13/cobra"
)

// NewCommand returns the device command of type cobra.Command
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "deviceservice",
		Short: "Device service command",
		Long: `Device service

Device Services (DS) are the edge connectors interacting with the devices 
or IoT objects that include, but are not limited to: appliances 
in your home, alarm systems, HVAC equipment, lighting, machines in any 
industry, irrigation systems, drones, traffic signals, automated transportation, 
and so forth.

Device services may service one or a number of devices, including sensors, 
actuators, and so forth, at one time. A “device” that a DS manages, 
could be something other than a simple single physical device and 
could be another gateway and all of that gateway’s devices, a device 
manager, or a device aggregator that acts as a device, or collection 
of devices, to EdgeX Foundry.

The Device Services layer’s microservices communicate with the devices, 
sensors, actuators, and other IoT objects through protocols native to the IoT object. 
The DS Layer converts the data produced and communicated by the IoT object, into a 
common EdgeX Foundry data structure, and sends that converted data into the Core Services 
layer, and to other microservices in other layers of EdgeX Foundry.`,
	}
	cmd.AddCommand(rm.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(add.NewCommand())
	cmd.AddCommand(update.NewCommand())
	return cmd
}
