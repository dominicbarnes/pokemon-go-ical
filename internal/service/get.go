package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/internal/ical"
)

const sourceURL = "https://raw.githubusercontent.com/bigfoott/ScrapedDuck/data/events.min.json"

func (s *Service) GetCalendar(ctx context.Context, id string) (string, error) {
	var cfg *CalendarConfig
	if err := s.DB.Get(ctx, id).ScanDoc(&cfg); err != nil {
		return "", fmt.Errorf("failed to get document: %w", err)
	}

	events, err := s.getEvents(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get events: %w", err)
	}

	tz, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return "", fmt.Errorf("timezone %s is invalid: %w", cfg.Timezone, err)
	}

	options := ical.GenerateICalOptions{
		Now:          time.Now(),
		TZ:           tz,
		IncludeTypes: cfg.IncludeEvents,
		ExcludeTypes: cfg.ExcludeEvents,
	}

	ics, err := ical.GenerateICal(events, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate ics: %w", err)
	}

	return ics.Serialize(), nil
}

func (s *Service) getEvents(ctx context.Context) ([]ical.LeekDuckEvent, error) {
	res, err := s.HTTP.Get(sourceURL)
	if err != nil {
		return nil, errors.New("failed to download events")
	}
	defer res.Body.Close()

	var ee []ical.LeekDuckEvent
	d := json.NewDecoder(res.Body)
	if err := d.Decode(&ee); err != nil {
		return nil, fmt.Errorf("failed to decode events from leek duck as JSON: %w", err)
	}
	return ee, nil
}
