package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/edgexfoundry-holding/edgex-cli/pkg/urlclient"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/spf13/viper"
)

func Get(url string, items interface{}) (err error) {
	ctx, _ := context.WithCancel(context.Background())
	printURL(url)
	data, err := clients.GetRequest(ctx, "", urlclient.NewA(url))
	if err != nil {
		return err
	}
	printResult(string(data))
	return json.Unmarshal(data, &items)
}

func Delete(url string) error {
	ctx, _ := context.WithCancel(context.Background())
	printURL(url)
	return clients.DeleteRequest(ctx, "", urlclient.NewA(url))
}


// DeleteItemByIdOrName deletes the given item
// The ID parameter can be either NAME or ID. We are doing this to allow the user
// enter either the name or the ID of an object to delete.
// First, we try ID. If successful, stop. If unsuccessful, try name.

//depricated
func DeleteItemByIdOrName(id string, pathID string, pathName string, url string) error {
	// Try to delete the object by Id
	err := Delete(url+pathID+id)
	if err != nil {
		// Try to delete the object by name
		if pathName == "" {
			return errors.New("Deleting by ID failed: " + url)
		}
		err = Delete(url+pathName+id)
	}
	return err
}

func printURL(url string) {
	if viper.GetBool("url") || viper.GetBool("verbose") {
		fmt.Printf("%s:%s \n",http.MethodGet, url)
	}
}

func printResult(data string) {
	if viper.GetBool("verbose") {
		fmt.Println(data)
	}
}