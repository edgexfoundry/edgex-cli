package client

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// GetAllItems returns a list of all Items in the DB
func GetAllItems(itemType string) []byte {

	host := viper.GetString("Host")

	resp, err := http.Get("http://" + host + ":48081/api/v1/" + itemType)
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

	host := viper.GetString("Host")

	// Create client
	client := &http.Client{}
	// call /Item/id/{id}
	req, err := http.NewRequest("DELETE", "http://"+host+":48081/api/v1/"+itemType+"/id/"+id, nil)
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

	// respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

}
