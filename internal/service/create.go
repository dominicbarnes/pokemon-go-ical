package service

import (
	"context"
	"fmt"
	"time"
)

type CreateCalendarInput struct {
	Timezone      string   `json:"timezone,omitempty"`
	IncludeEvents []string `json:"include_events,omitempty"`
	ExcludeEvents []string `json:"exclude_events,omitempty"`
}

func (s *Service) CreateCalendar(ctx context.Context, input *CreateCalendarInput) (*CalendarConfig, error) {
	id, rev, err := s.DB.CreateDoc(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	tz, err := time.LoadLocation(input.Timezone)
	if err != nil {
		return nil, fmt.Errorf("timezone %s is invalid: %w", input.Timezone, err)
	}

	return &CalendarConfig{
		ID:            id,
		Rev:           rev,
		Timezone:      tz.String(),
		IncludeEvents: input.IncludeEvents,
		ExcludeEvents: input.ExcludeEvents,
	}, nil
}
