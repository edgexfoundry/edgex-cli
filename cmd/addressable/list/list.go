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

package list

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type addressableList struct {
	rd []models.Addressable
}

// NewCommand returns the list device command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A list of all device services",
		Long:  `Return all device services sorted by id.`,
		Run: func(cmd *cobra.Command, args []string) {

			data, err := client.GetAllItems(
				"addressable",
				config.Conf.MetadataService.Port)

			if data == nil {
				return
			}

			if err != nil {
				fmt.Println(err)
				return
			}

			addressableList1 := addressableList{}

			errjson := json.Unmarshal(data, &addressableList1.rd)
			if errjson != nil {
				fmt.Println(errjson)
			}

			pw := viper.Get("writer").(io.Writer)
			w := new(tabwriter.Writer)
			w.Init(pw, 0, 8, 1, '\t', 0)
			_, err = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "ID", "Name", "Protocol")
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, addressable := range addressableList1.rd {
				fmt.Fprintf(w, "%s\t%s\t%v\t\n",
					addressable.Id,
					addressable.Name,
					addressable.Protocol,
				)
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = w.Flush()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	return cmd
}
