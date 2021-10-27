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
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:          "metrics",
		Short:        "Output the CPU/memory usage stats for all EdgeX core/support microservices",
		Long:         "Output the CPU/memory usage stats for all EdgeX core/support microservices",
		RunE:         handleMetrics,
		SilenceUsage: true,
	}
	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)
}

func handleMetrics(cmd *cobra.Command, args []string) error {
	var w *tabwriter.Writer
	services := getSelectedServices()

	if !json {
		w = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(w, "Service\tCpuBusyAvg\tMemAlloc\tMemFrees\tMemLiveObjects\tMemMallocs\tMemSys\tMemTotalAlloc")
		defer w.Flush()
	}

	for serviceName, service := range services {
		client := service.GetCommonClient()
		response, err := client.Metrics(context.Background())
		if err == nil {
			if json {
				result, err := jsonpkg.Marshal(response)
				if err != nil {
					return err
				}
				fmt.Println(string(result))
			} else {
				if err == nil {
					fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
						serviceName,
						response.Metrics.CpuBusyAvg,
						response.Metrics.MemAlloc,
						response.Metrics.MemFrees,
						response.Metrics.MemLiveObjects,
						response.Metrics.MemMallocs,
						response.Metrics.MemSys,
						response.Metrics.MemTotalAlloc)

				}
			}
		}
	}
	return nil
}
