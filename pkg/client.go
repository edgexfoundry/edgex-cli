package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

var client = &http.Client{}

// GetAllItemsDepricated returns a list of all Items in the DB
func GetAllItems(url string) ([]byte, error) {

	// Get URL and VERBOSE from viper
	urlFlag := viper.GetBool("url")
	verboseFlag := viper.GetBool("verbose")

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

func DeleteItem(url string) ([]byte, error) {

	urlFlag := viper.GetBool("url")
	verboseFlag := viper.GetBool("verbose")

	if urlFlag {
		fmt.Println("DELETE: " + url)
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

// DeleteItemByIdOrName deletes the given item
// The ID parameter can be either NAME or ID. We are doing this to allow the user
// enter either the name or the ID of an object to delete.
// First, we try ID. If successful, stop. If unsuccessful, try name.

//depricated
func DeleteItemByIdOrName(id string, pathID string, pathName string, url string) ([]byte, error) {
	// Try to delete the object by Id
	respBody, err := DeleteItem(url+pathID+id)
	if string(respBody) == SUCCESSFUL_DELETE {
		// deleting with ID worked
		return respBody, err
	}

	// Try to delete the object by name
	if pathName == "" {
		return nil, errors.New("Deleting by ID failed: " + url)
	}
	respBody, err = DeleteItem(url+pathName+id)
	return respBody, err
}
