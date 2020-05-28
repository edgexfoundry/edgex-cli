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
	"github.com/spf13/viper"
	"html/template"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

var device string

const cmdTempl = "Id\tName\tPath\t\n" +
	"{{range $i, $c := .Commands}}" +
	"{{$c.Id}}\t{{$c.Name}}\t{{$c.Get.URL}}\t\n" +
	"{{end}}"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of device supported commands",
		Long:  "Return a list of all device `commands`",
		RunE:  listHandler,
	}
	cmd.Flags().StringVarP(&device, "device", "d", "", "List commands associated with device specified by name")
	cmd.MarkFlagRequired("device")
	return cmd
}

func listHandler(cmd *cobra.Command, args []string) (err error) {
	url := config.Conf.Clients["Command"].Url() + clients.ApiDeviceRoute + "/name/" + device
	var response models.CommandResponse
	err = request.Get(url, &response)
	if err != nil {
		return
	}
	pw := viper.Get("writer").(io.WriteCloser)
	w := new(tabwriter.Writer)
	w.Init(pw, 0, 8, 1, '\t', 0)
	tmpl, err := template.New("commandList").Parse(cmdTempl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, response)
	if err != nil {
		return err
	}
	w.Flush()
	return
}
