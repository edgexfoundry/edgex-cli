// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"

	"github.com/edgexfoundry/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/spf13/cobra"
)

// NewCommand return add profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [FILE]...",
		Short: "Add device profiles",
		Long:  `Upload the given YAML files to core-metadata for device profile creation.`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, fname := range args {
				dp, err := parseYaml(fname)
				if err != nil {
					fmt.Println("Error: ", err.Error())
				}
				addDeviceProfile(dp)
			}
		},
	}
	return cmd
}

func addDeviceProfile(dp *models.DeviceProfile) {
	url := config.Conf.Clients["Metadata"].Url()
	mdc := metadata.NewDeviceProfileClient(
		local.New(url + clients.ApiDeviceProfileRoute),
	)

	dpId, err := mdc.Add(context.Background(), dp)
	if err != nil {
		fmt.Printf("Failed to create Device Profile `%s` because of error: %s\n", dp.Name, err)
	} else {
		fmt.Printf("Device Profile successfully created: %s\n", dpId)
	}
}

func parseYaml(fname string) (*models.DeviceProfile, error) {
	var dp = &models.DeviceProfile{}
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, dp)
	if err != nil {
		return nil, err
	}

	return dp, nil
}
