package main

import (
	"time"
)

func parseLeekDuckEventTime(raw string, tz *time.Location) (*time.Time, error) {
	// Leek Duck prints events without a timezone suffix to indicate that they are
	// expected to be in the user's local time, rather than a predefined zone.
	localTime, err := time.ParseInLocation("2006-01-02T15:04:05", string(raw), tz)
	if err == nil {
		return &localTime, nil
	}
	// ignore errors from attempting to parse as local time

	parsedTime, err := time.Parse("2006-01-02T15:04:05.999-0700", string(raw))
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}
