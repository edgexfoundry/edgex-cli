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
	if viper.GetBool("url") {
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

func DeleteByIds(i interface{}, ids []string) error {
	methodVal := reflect.ValueOf(i).MethodByName("Delete")
	if !methodVal.IsValid() {
		return fmt.Errorf("unsupported method: %s", "Delete")
	} else if methodVal.Type().NumIn() != 2 {
		return fmt.Errorf("client method has %q input parameters, want 2", methodVal.Type().NumIn())
	} else if methodVal.Type().In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return fmt.Errorf("client method's first input parameter is %q, want `context.Context`", methodVal.Type().In(0))
	} else if methodVal.Type().In(1) != reflect.TypeOf((*string)(nil)).Elem() {
		return fmt.Errorf("client method's first input parameter is %q, want `string`", methodVal.Type().In(1))
	} else if methodVal.Type().NumOut() != 1 {
		return fmt.Errorf("client method has %q output parameters, want 1", methodVal.Type().NumOut())
	}

	for _, id := range ids {
		out := methodVal.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(id)})
		err := out[0]
		if !err.IsNil() {
			fmt.Printf("Error: %s", err.Interface().(error))
		} else {
			fmt.Printf("Removed: %s \n", id)
		}
	}
	return nil
}
