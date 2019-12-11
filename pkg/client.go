package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

var client = &http.Client{}

func buildURL(itemType string, path string, port string) string {
	host := viper.GetString("host")
	url := "http://" + host + ":" + port + "/api/v1/" + path + itemType
	return url
}

// GetAllItems returns a list of all Items in the DB
func GetAllItems(itemType string, port string) ([]byte, error) {

	// Get URL and VERBOSE from viper
	urlFlag := viper.GetBool("url")
	verboseFlag := viper.GetBool("verbose")

	url := buildURL(itemType, "", port)

	if urlFlag {
		fmt.Println("GET: " + url)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Need to read body from response
	body, errBody := ioutil.ReadAll(resp.Body)

	// If verbose is enabled, print Header and Body + errors if any
	if verboseFlag {
		fmt.Println("HEADER")
		for k, v := range resp.Header {
			fmt.Printf("%v : %v\n", k, v)
		}

		if errBody != nil {
			fmt.Println(errBody)
		} else {
			fmt.Println("BODY")
			fmt.Println(string(body))
		}
		return nil, nil
	}

	return body, errBody
}

// DeleteItemByID deletes an item by ID
func DeleteItemByID(id string, pathID string, port string) ([]byte, error) {

	urlFlag := viper.GetBool("url")
	verboseFlag := viper.GetBool("verbose")
	url := buildURL(id, pathID, port)

	if urlFlag {
		fmt.Println("GET: " + url)
	}

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

	body, errBody := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// If verbose is enabled, print Header and Body + errors if any
	if verboseFlag {
		fmt.Println("HEADER")
		for k, v := range resp.Header {
			fmt.Printf("%v : %v\n", k, v)
		}

		if errBody != nil {
			fmt.Println(errBody)
		} else {
			fmt.Println("BODY")
			fmt.Println(string(body))
		}
	}

	return body, nil
}

// DeleteItemByName deletes the item by name
func DeleteItemByName(id string, pathName string, port string) ([]byte, error) {

	urlFlag := viper.GetBool("url")
	verboseFlag := viper.GetBool("verbose")
	url := buildURL(id, pathName, port)

	if urlFlag {
		fmt.Println("GET: " + url)
	}

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

	respBody, errBody := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// If verbose is enabled, print Header and Body + errors if any
	if verboseFlag {
		fmt.Println("HEADER")
		for k, v := range resp.Header {
			fmt.Printf("%v : %v\n", k, v)
		}

		if errBody != nil {
			fmt.Println(errBody)
		} else {
			fmt.Println("BODY")
			fmt.Println(string(respBody))
		}
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
	urlFlag := viper.GetBool("url")
	respBody, err := DeleteItemByID(id, pathID, port)
	if urlFlag {
		fmt.Println("DELETE: " + url)
	}
	if string(respBody) == SUCCESSFUL_DELETE {
		// deleting with ID worked
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
