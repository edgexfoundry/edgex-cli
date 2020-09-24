// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package rm

import (
	"errors"
	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

var name string

// NewCommand returns the rm device service command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "rm [--name|<id>]",
		Example: "deviceservice rm <id> \n" +
			"deviceservice rm --name <name>",
		Short: "Removes device service by name or ID",
		Long:  `Removes a device service from the core-metadata DB.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 && name == "" {
				return errors.New("no device service id/name provided")
			}
			url := config.Conf.Clients["Metadata"].Url() + clients.ApiDeviceServiceRoute
			url, deletedBy := constructUrl(url, args)
			return request.DeletePrt(url, deletedBy)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Delete Device Service by given name")
	return cmd
}

func constructUrl(url string, args []string) (string, string) {
	if name != "" {
		url = url + config.PathName + name
		return url, name
	}
	url = url + config.PathId + args[0]
	return url, args[0]
}
