package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalendarCreateInput struct {
	Timezone      string   `json:"timezone,omitempty"`
	IncludeEvents []string `json:"include_events,omitempty"`
	ExcludeEvents []string `json:"exclude_events,omitempty"`
}

type CalendarCreateOutput struct {
	ID            string   `json:"id,omitempty"`
	Timezone      string   `json:"timezone,omitempty"`
	IncludeEvents []string `json:"include_events,omitempty"`
	ExcludeEvents []string `json:"exclude_events,omitempty"`
}

func (s *Service) CalendarCreate(c echo.Context) error {
	ctx := c.Request().Context()

	var input CalendarCreateInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	id, _, err := s.KivikDB.CreateDoc(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	tz, err := validateTZ(input.Timezone)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/calendars/%s", id))
	return c.JSON(http.StatusCreated, CalendarCreateOutput{
		ID:            id,
		Timezone:      tz.String(),
		IncludeEvents: input.IncludeEvents,
		ExcludeEvents: input.ExcludeEvents,
	})
}
