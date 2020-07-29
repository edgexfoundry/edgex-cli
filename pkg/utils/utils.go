package utils

import (
	"fmt"
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

const TimeUnitDescriptions = "ms - milliseconds\n" +
	"s - seconds\n" +
	"m - minutes\n" +
	"h - hours\n" +
	"d - days\n"

//empty struct has width of zero. It occupies zero bytes of storage
var TimeUnitsMap = map[string]struct{}{"ms": {}, "s": {}, "h": {}, "d": {}, "m": {}}

func ConvertAgeToMillisecond(unit string, age int64) int64 {
	var ageMilliseconds int64
	switch unit {
	case "ms":
		ageMilliseconds = age
	case "s":
		ageMilliseconds = age * 1000
	case "m":
		ageMilliseconds = age * 60 * 1000
	case "h":
		ageMilliseconds = age * 60 * 60 * 1000
	case "d":
		ageMilliseconds = age * 24 * 60 * 60 * 1000
	}
	return ageMilliseconds
}

