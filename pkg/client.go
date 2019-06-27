package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetAllItems returns a list of all Items in the DB
func GetAllItems(itemType string) []byte {

	resp, err := http.Get("http://localhost:48081/api/v1/" + itemType)
	if err != nil {
		// handle error
		fmt.Println("An error occurred")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	return data
}

// DeleteItem deletes the given item
func DeleteItem(id string, itemType string) {

	// Create client
	client := &http.Client{}
	// call /Item/id/{id}
	req, err := http.NewRequest("DELETE", "http://localhost:48081/api/v1/"+itemType+"/"+id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

}
