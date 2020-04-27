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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"

	"github.com/edgexfoundry-holding/edgex-cli/cmd/addressable"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/db"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/device"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/deviceservice"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/event"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/interval"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/notification"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/profile"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/reading"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/status"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/subscription"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/version"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/pager"
)

// NewCommand returns rootCmd which represents the base command when called without any subcommands
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set flags
			noPager, err := cmd.Flags().GetBool("no-pager")
			if err != nil {
				fmt.Println("couldn't get no-pager flag")
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				fmt.Println("couldn't get verbose flag")
			}
			viper.Set("verbose", verbose)
			if verbose {
				noPager = true
			}

			url, err := cmd.Flags().GetBool("url")
			if err != nil {
				fmt.Println("couldn't get url flag")
			}
			viper.Set("url", url)

			configFile, err := cmd.Flags().GetString("config-file")
			if err != nil {
				fmt.Println("couldn't get config-file flag")
			}
			viper.Set("config-file", configFile)

			viper.Set("writer", os.Stdout)
			if !noPager {
				w, err := pager.NewWriter()
				if err == nil {
					viper.Set("writer", w)
					viper.Set("writerShouldClose", true) // This flag prevents us from calling close on stdout
				}
			}

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			shouldClose := viper.GetBool("writerShouldClose")
			if shouldClose {
				pw := viper.Get("writer").(io.Closer)
				if pw != os.Stdout {
					err := pw.Close()
					if err != nil {
						_ = fmt.Errorf(err.Error())
					}
				}
			}
		},
		Use:   "edgex-cli",
		Short: "EdgeX command line interface",
		Long: "\n" +

    " ______     _              __   __            _____  _       _____        \n" +
	"|  ____|   | |             \\ \\ / /           / ____|| |     |_   _|     \n" +
	"| |__    __| |  __ _   ___  \\ V /   ______  | |     | |       | |        \n" +
	"|  __|  / _` | / _` | / _ \\  > <   |______| | |     | |       | |        \n" +
	"| |____| (_| || (_| ||  __/ / . \\           | |____ | |____  _| |_       \n" +
	"|______|\\__,_| \\__, | \\___|/_/ \\_\\           \\_____||______||_____| \n" +
	"		__/ |                                                             \n" +
	"	       |___/                                                              \n" +


	`
EdgeX CLI version: ` + version.Version +
			`
https://www.edgexfoundry.org/
	`,
	}

	// Add all subcommands below:
	cmd.AddCommand(device.NewCommand())
	cmd.AddCommand(deviceservice.NewCommand())
	cmd.AddCommand(profile.NewCommand())
	cmd.AddCommand(event.NewCommand())
	cmd.AddCommand(reading.NewCommand())
	cmd.AddCommand(status.NewCommand())
	cmd.AddCommand(db.NewCommand())
	cmd.AddCommand(addressable.NewCommand())
	// --- Support Services Commands ---
	cmd.AddCommand(notification.NewCommand())
	cmd.AddCommand(subscription.NewCommand())
	cmd.AddCommand(interval.NewCommand())
	cmd.AddCommand(version.NewCommand())

	// global flags
	Verbose := false
	URL := false
	NoPager := false
	Configfile := ""
	// get flags values
	cmd.PersistentFlags().BoolVarP(&URL, "url", "u", false, "Print URL(s) used by the entered command.")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print entire HTTP response.")
	cmd.PersistentFlags().BoolVarP(&NoPager, "no-pager", "", false, "Do not pipe output into a pager.")
	cmd.PersistentFlags().StringVar(&Configfile, "config-file", "", "configuration file")
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//var ConfigFile string
	//flag.StringVar(&ConfigFile, "config-file", config.DefaultConfigFile, "Specify local configuration directory")
	var env = config.NewViperEnv()
	// should we really be reading config file later? when does the cmd persistent flags get read .. after that to pull in any config-file flag argument?
	env.SetConfigFile(config.DefaultConfigFile)
	if err := config.LoadConfig(env); err != nil {
		os.Exit(1)
	}
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
