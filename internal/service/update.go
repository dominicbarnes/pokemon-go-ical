package service

import (
	"context"
	"fmt"
	"time"
)

type UpdateCalendarInput struct {
	ID            string    `json:"id"`
	Timezone      *string   `json:"timezone,omitempty"`
	IncludeEvents *[]string `json:"include_events,omitempty"`
	ExcludeEvents *[]string `json:"exclude_events,omitempty"`
}

func (s *Service) UpdateCalendar(ctx context.Context, input *UpdateCalendarInput) (*CalendarConfig, error) {
	var doc *CalendarConfig
	if err := s.DB.Get(ctx, input.ID).ScanDoc(&doc); err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if input.Timezone != nil {
		tz, err := time.LoadLocation(*input.Timezone)
		if err != nil {
			return nil, fmt.Errorf("timezone %s is invalid: %w", input.Timezone, err)
		}

		doc.Timezone = tz.String()
	}

	if input.IncludeEvents != nil {
		doc.IncludeEvents = *input.IncludeEvents
	}

	if input.ExcludeEvents != nil {
		doc.ExcludeEvents = *input.ExcludeEvents
	}

	rev, err := s.DB.Put(ctx, doc.ID, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to put document: %w", err)
	}
	doc.Rev = rev

	return doc, nil
}
