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
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/edgexfoundry/edgex-cli/cmd/db"
	"github.com/edgexfoundry/edgex-cli/cmd/device"
	"github.com/edgexfoundry/edgex-cli/cmd/deviceservice"
	"github.com/edgexfoundry/edgex-cli/cmd/event"
	"github.com/edgexfoundry/edgex-cli/cmd/interval"
	"github.com/edgexfoundry/edgex-cli/cmd/notification"
	"github.com/edgexfoundry/edgex-cli/cmd/profile"
	"github.com/edgexfoundry/edgex-cli/cmd/reading"
	"github.com/edgexfoundry/edgex-cli/cmd/status"
	"github.com/edgexfoundry/edgex-cli/cmd/subscription"
	"github.com/edgexfoundry/edgex-cli/cmd/version"
	"github.com/edgexfoundry/edgex-cli/config"
)

// NewCommand returns rootCmd which represents the base command when called without any subcommands
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edgex",
		Short: "EdgeX command line interface",
		Long: " _____    _           __  __  _____                     _            \n" +
			"| ____|__| | __ _  ___\\ \\/ / |  ___|__  _   _ _ __   __| |_ __ _   _ \n" +
			"|  _| / _` |/ _` |/ _ \\\\  /  | |_ / _ \\| | | | '_ \\ / _` | '__| | | |\n" +
			"| |__| (_| | (_| |  __//  \\  |  _| (_) | |_| | | | | (_| | |  | |_| |\n" +
			"|_____\\__,_|\\__, |\\___/_/\\_\\ |_|  \\___/ \\__,_|_| |_|\\__,_|_|   \\__, |\n" +
			"            |___/                                              |___/ \n" +
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
	// --- Support Services Commands ---
	cmd.AddCommand(notification.NewCommand())
	cmd.AddCommand(subscription.NewCommand())
	cmd.AddCommand(interval.NewCommand())
	cmd.AddCommand(version.NewCommand())

	// global flags
	Verbose := false
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print URL(s) used by the entered command.")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// set default config
	setConfig()
	if err := NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setConfig() {
	// Config file path
	configFilePath := os.Getenv("HOME") + "/.edgex-cli/config.yaml"

	// Set config file
	viper.SetConfigFile(configFilePath)

	var configuration config.Configuration

	// checking if file already exists
	if !exists(configFilePath) {
		var defaultConfig = []byte(`
host: "localhost"
security:
  enabled: false
  token: "empty"
ports:
  CoreData: "48080" 
  CoreMetadata: "48081"
  CoreCommand: "48082" 
  Notifications: "48060"  
  Logging: "48061" 
  Scheduling: "48085"
  RulesEngine: "48075"
  ClientRegistration: "48071"
  SystemManagement: "48090"
`)

		f, err := os.Create(configFilePath)
		if err != nil {
			log.Fatalf("Error creating config file, %s", err)
		}
		defer f.Close()

		_, err = f.Write(defaultConfig)
		if err != nil {
			log.Fatalf("Error write config file, %s", err)
		}

	}

	// Reading from file that was already existing or newly created
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		// log.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	viper.WriteConfig()
	viper.SafeWriteConfig()
}

// Helper function to check whether file exists
func exists(configPath string) bool {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
