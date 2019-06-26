// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCommand return add profile command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [FILE]...",
		Short: "Add device profiles",
		Long:  `Upload the given YAML files to core-metadata for device profile creation.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add Device Profiles:")
			client := &http.Client{}
			for _, val := range args {
				fmt.Print(val, "... ")
				// Open file
				data, err := os.Open(val)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer data.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", filepath.Base(val))
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = io.Copy(part, data)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = writer.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
				// Create request
				req, err := http.NewRequest("POST", "http://localhost:48081/api/v1/deviceprofile/uploadfile", body)
				if err != nil {
					fmt.Println(err)
					return
				}
				req.Header.Add("Content-Type", writer.FormDataContentType())

				// Fetch Request
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()

				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Display Results
				if resp.StatusCode == 200 {
					fmt.Println("OK, new profile ID is", string(respBody))
				} else {
					fmt.Print(resp.Status, ": ", string(respBody))
				}
			}
		},
	}
	return cmd
}
