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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:          "metrics",
		Short:        "Outputs the CPU/memory usage stats for all EdgeX core/support microservices",
		Long:         ``,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			err = showMetrics(cmd)
			return err
		}}

	rootCmd.AddCommand(cmd)
	addStandardFlags(cmd)

}

func showMetrics(cmd *cobra.Command) error {
	services := getSelectedServices()

	if json || verbose {
		for serviceName, service := range services {
			jsonValue, url, err := service.GetMetricsJSON()
			if err != nil {
				return err
			} else {
				if json {
					cmd.Println(jsonValue)
				} else {
					cmd.Printf("%s: %s: %s\n", serviceName, url, jsonValue)
				}
			}
		}
	} else {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(w, "Service\tCpuBusyAvg\tMemAlloc\tMemFrees\tMemLiveObjects\tMemMallocs\tMemSys\tMemTotalAlloc")
		for serviceName, service := range services {
			metrics, err := service.GetMetrics()
			if err == nil {
				fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", serviceName, metrics.CpuBusyAvg, metrics.MemAlloc, metrics.MemFrees,
					metrics.MemLiveObjects, metrics.MemMallocs, metrics.MemSys, metrics.MemTotalAlloc)

			}
		}
		w.Flush()

	}
	return nil
}
