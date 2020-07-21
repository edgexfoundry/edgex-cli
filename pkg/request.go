package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"

	"github.com/spf13/viper"
)

func Get(url string, items interface{}) (err error) {
	printURL(url, http.MethodGet)
	resp, err := clients.GetRequest(context.Background(), "", local.New(url))
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, &items)
}
func Delete(url string) error {
	printURL(url, http.MethodDelete)
	return clients.DeleteRequest(context.Background(), "", local.New(url))
}
func DeletePrt(url string, deletedBy string) error {
	err := Delete(url)
	if err == nil && deletedBy != "" {
		fmt.Printf("Removed: %s \n", deletedBy)
		return nil
	}
	return err
}

func Post(url string, item interface{}) {
	printURL(url, http.MethodPost)
	printData(item)
	resp, err := clients.PostJSONRequest(context.Background(), "", item, local.New(url))
	name := getType(item)
	if err != nil {
		fmt.Printf("Failed to create %s because of error: %s", name, err)
	} else {
		fmt.Printf("%s successfully created: %s ", name, resp)
	}
}

func printURL(url string, method string) {
	if viper.GetBool("url") || viper.GetBool("verbose") {
		fmt.Printf("> %s: %s \n", method, url)
	}
}

func printData(item interface{}) {
	if viper.GetBool("verbose") {
		body := reflect.ValueOf(item).MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)
		fmt.Printf("> Request data: %s\n", body)
	}
}

func getType(item interface{}) string {
	if t := reflect.TypeOf(item); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
