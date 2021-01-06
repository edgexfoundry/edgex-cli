/*******************************************************************************
 * Copyright 2020 VMWare.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package list

import (
	"context"
	"html/template"
	"strings"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"
	"github.com/edgexfoundry/edgex-cli/pkg/formatters"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var device string

const cmdsTempl = "Name\tDevice Id\tDevice Name\tMethods\tURL\t\n" +
	"{{range .}}" +
	"{{ $deviceId :=.Id}}" +
	"{{ $deviceName :=.Name}}" +
	"{{range $i, $c := .Commands}}" +
	"{{$c.Name}}\t{{$deviceId}}\t{{$deviceName}}\t{{supportedMethods $c}}\t{{if $c.Get.URL}}{{$c.Get.URL}}{{else}}{{$c.Put.URL}}{{end}}\t\n" +
	"{{end}}" +
	"{{end}}"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of device supported commands",
		Long:  "Return a list of all device `commands`",
		RunE:  listHandler,
	}
	cmd.Flags().StringVarP(&device, "device", "d", "", "List commands associated with device specified by name")
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	var responses []models.CommandResponse
	if device != "" {
		responses, err = getCommandsByDeviceName(cmd.Context(), device)
	} else {
		responses, err = getCommands(cmd.Context())
	}
	if err != nil {
		return
	}
	formatter := formatters.NewFormatter(cmdsTempl, template.FuncMap{"supportedMethods": supportedMethods})
	err = formatter.Write(responses)
	return
}

func getCommands(ctx context.Context) ([]models.CommandResponse, error) {
	var responses []models.CommandResponse
	url := config.Conf.Clients["Command"].Url() + clients.ApiDeviceRoute
	err := request.Get(ctx, url, &responses)
	return responses, err
}

func getCommandsByDeviceName(ctx context.Context, d string) ([]models.CommandResponse, error) {
	var response models.CommandResponse
	url := config.Conf.Clients["Command"].Url() + clients.ApiDeviceRoute + "/name/" + device
	err := request.Get(ctx, url, &response)
	responses := []models.CommandResponse{response}
	return responses, err
}

func supportedMethods(cmd models.Command) (methods string) {
	if cmd.Get.Path != "" {
		methods = "Get,"
	}
	if cmd.Put.Path != "" {
		methods = methods + " Put"
	}
	methods = strings.TrimSpace(methods)
	return strings.TrimSuffix(methods, ",")
}
