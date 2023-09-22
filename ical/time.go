package main

import (
	"fmt"
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

	for _, format := range []string{"2006-01-02T15:04:05.999-0700", time.RFC3339Nano} {
		parsedTime, err := time.Parse(format, string(raw))
		if err == nil {
			return &parsedTime, nil
		}
	}

	return nil, fmt.Errorf("invalid event time: %s", raw)
}
