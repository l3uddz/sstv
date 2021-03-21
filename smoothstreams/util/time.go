package util

import (
	"fmt"
	"time"
)

var (
	locSmoothStreams *time.Location
)

func init() {
	// time locations
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("load time location: %v", err))
	}
	locSmoothStreams = loc
}

func CurrentXmlTvTime() string {
	t := time.Now().UTC().In(locSmoothStreams)
	return t.Format("20060102150405 -0700")
}

func TimeStringToXmlTvTime(ts string) (string, error) {
	// parse time from smoothstreams (comes in EST/EDT - New York)
	t, err := time.ParseInLocation("2006-01-02 15:04:05 EST", ts+" EST", locSmoothStreams)
	if err != nil {
		return "", fmt.Errorf("parse in location: %v: %w", ts, err)
	}

	return t.UTC().Format("20060102150405 -0700"), nil
}
