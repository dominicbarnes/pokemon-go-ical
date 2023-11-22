package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalendarUpdateInput struct {
	ID            string    `param:"id"`
	Timezone      *string   `json:"timezone,omitempty"`
	IncludeEvents *[]string `json:"include_events,omitempty"`
	ExcludeEvents *[]string `json:"exclude_events,omitempty"`
}

type CalendarUpdateOutput struct {
	ID            string   `json:"id,omitempty"`
	Timezone      string   `json:"timezone,omitempty"`
	IncludeEvents []string `json:"include_events,omitempty"`
	ExcludeEvents []string `json:"exclude_events,omitempty"`
}

func (s *Service) CalendarUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	var input CalendarUpdateInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	var doc *CalendarConfig
	if err := s.KivikDB.Get(ctx, input.ID).ScanDoc(&doc); err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	if input.Timezone != nil {
		tz, err := validateTZ(*input.Timezone)
		if err != nil {
			return err
		}

		doc.Timezone = tz.String()
	}

	if input.IncludeEvents != nil {
		doc.IncludeEvents = *input.IncludeEvents
	}

	if input.ExcludeEvents != nil {
		doc.ExcludeEvents = *input.ExcludeEvents
	}

	rev, err := s.KivikDB.Put(ctx, doc.ID, doc)
	if err != nil {
		return fmt.Errorf("failed to put document: %w", err)
	}
	doc.Rev = rev

	return c.JSON(http.StatusOK, CalendarUpdateOutput{
		ID:            doc.ID,
		Timezone:      doc.Timezone,
		IncludeEvents: doc.IncludeEvents,
		ExcludeEvents: doc.ExcludeEvents,
	})
}
