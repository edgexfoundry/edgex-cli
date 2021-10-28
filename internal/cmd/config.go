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
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:          "config",
		Short:        "Return the current configuration of all EdgeX core/support microservices",
		Long:         "Return the current configuration of all EdgeX core/support microservices",
		RunE:         handleConfig,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)

}

func handleConfig(cmd *cobra.Command, args []string) error {
	services := getSelectedServices()

	for serviceName, service := range services {
		client := service.GetCommonClient()
		response, err := client.Configuration(context.Background())
		if err == nil {
			jsonresult, xerr := jsonpkg.Marshal(response)
			if xerr != nil {
				return xerr
			}

			if json {
				fmt.Println(string(jsonresult))
			} else {
				fmt.Println(serviceName + ":")
				var result map[string]interface{}
				jsonpkg.Unmarshal([]byte(jsonresult), &result)
				b, err := jsonpkg.MarshalIndent(result["config"], "", "    ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			}
		}
	}
	return nil
}
