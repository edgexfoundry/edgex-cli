package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// humanDuration return the duration since date. Taken from https://github.com/docker/go-units/blob/master/duration.go
func humanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds == 1 {
		return "1 second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours() + 0.5); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*2 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
}

func DisplayDuration(tt int64) string {
	zeroDisplay := ""
	if tt == 0 {
		return zeroDisplay
	} else {
		ttTime := time.Unix(tt/1000, 0)
		return humanDuration(time.Since(ttTime))
	}
}

func ListHelper(url string, readings interface{}) (err error) {
	resp, err := http.Get(url)
	if err != nil  {
		// handle error
		fmt.Println("An error occurred. Is EdgeX running?")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println(err1)
		return err1
	} else if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("Server error: %d %s \n", resp.StatusCode, string(data))
	}

	err = json.Unmarshal(data, &readings)
	if err != nil {
		if string(data) == "Error, exceeded the max limit as defined in config" {
			fmt.Println("The number of readings to be returned exceeds the MaxResultCount limit defined in configuration.toml")
		}
		fmt.Println(err)
		return
	}
	return
}
