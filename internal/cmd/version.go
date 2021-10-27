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
		Use:          "version",
		Short:        "Output the current version of EdgeX CLI and EdgeX microservices",
		Long:         "Output the current version of EdgeX CLI and EdgeX microservices",
		RunE:         handleVersion,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)

}

func handleVersion(cmd *cobra.Command, args []string) error {
	services := getSelectedServices()

	for serviceName, service := range services {
		client := service.GetCommonClient()
		response, err := client.Version(context.Background())
		if err == nil {
			if json {
				result, err := jsonpkg.Marshal(response)
				if err != nil {
					return err
				}
				fmt.Println(string(result))
			} else {

				fmt.Println(serviceName + ": " + response.Version)
			}
		}
	}

	return nil
}
