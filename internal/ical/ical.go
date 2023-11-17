package ical

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
	cal.SetXWRCalName("Pokémon GO Events")
	cal.SetDescription("Powered by ScrapedDuck and LeekDuck.com")
	cal.SetXWRCalDesc("Powered by ScrapedDuck and LeekDuck.com")
	cal.SetXWRTimezone(options.TZ.String())
	cal.SetRefreshInterval("PT12H")
	cal.SetXPublishedTTL("PT12H")

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
		e.SetDtStampTime(options.Now)
		e.SetStartAt(*startAt)
		e.SetEndAt(*endAt)
		e.SetSummary(event.Title())
		e.SetDescription(event.Description())
		e.SetURL(event.Link)
		e.SetProperty("IMAGE;VALUE=URI", event.Image)

		a := e.AddAlarm()
		a.SetAction(ical.ActionDisplay)
		a.SetTrigger("-PT15M")
	}

	return cal, nil
}
