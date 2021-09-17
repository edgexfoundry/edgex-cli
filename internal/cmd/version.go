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

	"github.com/edgexfoundry/edgex-cli"
	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Outputs the current versions of EdgeX CLI and EdgeX microservices",
		Long:  ``,

		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if !json {
				cmd.Println("EdgeX CLI version: ", edgex.BuildVersion)
			}
			err = showVersion(cmd)

			return err
		}}

	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)

}

func showVersion(cmd *cobra.Command) error {
	services := getSelectedServices()

	for serviceName, service := range services {
		jsonValue, url, err := service.GetVersionJSON()
		if err != nil {
			if json {
				return err
			} else if verbose {
				cmd.Printf("%s: %s: %s\n", serviceName, url, err.Error())
			}
		} else {
			if json {
				cmd.Println(jsonValue)
			} else if verbose {
				cmd.Printf("%s: %s: %s\n", serviceName, url, jsonValue)
			} else {
				var result map[string]interface{}
				jsonpkg.Unmarshal([]byte(jsonValue), &result)
				cmd.Println(serviceName + ": " + result["version"].(string))
			}
		}
	}

	return nil
}
