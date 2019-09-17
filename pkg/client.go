package client

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{}

func buildURL(itemType string, path string, port string) string {
	host := viper.GetString("host")
	url := "http://" + host + ":" + port + "/api/v1/" + path + itemType
	return url
}

// GetAllItems returns a list of all Items in the DB
func GetAllItems(itemType string, port string, verbose bool) ([]byte, error) {

	url := buildURL(itemType, "", port)

	if verbose {
		fmt.Println("GET: " + url)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func DeleteItemByID(id string, pathID string, port string) ([]byte, error) {

	url := buildURL(id, pathID, port)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func DeleteItemByName(id string, pathName string, port string) ([]byte, error) {

	url := buildURL(id, pathName, port)
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// DeleteItem deletes the given item
// The ID parameter can be either NAME or ID. We are doing this to allow the user
// enter either the name or the ID of an object to delete.
// First, we try ID. If successful, stop. If unsuccessful, try name.

func DeleteItem(id string, pathID string, pathName string, port string, verbose bool) ([]byte, error) {
	// Try ID first
	url := buildURL(id, pathID, port)
	respBody, err := DeleteItemByID(id, pathID, port)

	if string(respBody) == SUCCESSFUL_DELETE {
		// deleting with ID worked
		if verbose {
			fmt.Println("DELETE: " + url)
		}
		return respBody, err
	}

	if pathName == "" {
		return nil, errors.New("Deleting by ID failed: " + url)
	}

	respBody, err = DeleteItemByName(id, pathName, port)
	return respBody, err
}

// GetVersion returns the version of a service given its port
func GetVersion(port string) []byte {
	host := viper.GetString("host")

	url := "http://" + host + ":" + port + "/version"

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Println("An error occurred")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	return data
}
