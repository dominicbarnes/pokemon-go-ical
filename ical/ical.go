package main

import (
	"fmt"
	"time"

	ical "github.com/arran4/golang-ical"
	mapset "github.com/deckarep/golang-set/v2"
)

type GenerateICalOptions struct {
	Now          time.Time
	TZ           *time.Location
	IncludeTypes []string
	ExcludeTypes []string
}

func GenerateICal(events []LeekDuckEvent, options GenerateICalOptions) (*ical.Calendar, error) {
	cal := ical.NewCalendar()
	cal.SetName("Pokémon GO Events")
	cal.SetDescription("Powered by ScrapedDuck and LeekDuck.com")
	cal.SetXWRTimezone(options.TZ.String())

	include := mapset.NewSet(options.IncludeTypes...)
	exclude := mapset.NewSet(options.ExcludeTypes...)

	for _, event := range events {
		if include.Cardinality() > 0 {
			if !include.Contains(event.Type) {
				continue
			}
		} else if exclude.Cardinality() > 0 {
			if exclude.Contains(event.Type) {
				continue
			}
		}

		startAt, err := parseLeekDuckEventTime(event.Start, options.TZ)
		if err != nil {
			return nil, fmt.Errorf("failed to parse event start %s: %w", event.Start, err)
		}

		endAt, err := parseLeekDuckEventTime(event.End, options.TZ)
		if err != nil {
			return nil, fmt.Errorf("failed to parse event start %s: %w", event.Start, err)
		}

		e := cal.AddEvent(event.ID)
		e.SetStartAt(*startAt)
		e.SetEndAt(*endAt)
		e.SetSummary(event.Title())
		e.SetURL(event.Link)
	}

	return cal, nil
}
