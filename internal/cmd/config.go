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
	jsonpkg "encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "config",
		Short: "Returns the current configuration of all EdgeX core/support microservices",
		Long:  ``,

		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = showConfig(cmd)
			return err
		}}

	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)

}

func showConfig(cmd *cobra.Command) error {
	services := getSelectedServices()

	for serviceName, service := range services {
		jsonValue, url, err := service.GetConfigJSON()
		if err != nil {
			if json {
				return err
			} else if verbose {
				fmt.Printf("%s: %s: %s\n", serviceName, url, err.Error())
			}
		} else {
			if json {
				fmt.Println(jsonValue)
			} else if verbose {
				fmt.Printf("%s: %s: %s\n", serviceName, url, jsonValue)
			} else {
				fmt.Println(serviceName + ":")
				var result map[string]interface{}
				jsonpkg.Unmarshal([]byte(jsonValue), &result)
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
