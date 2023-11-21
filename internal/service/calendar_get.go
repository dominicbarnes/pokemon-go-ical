package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kivik/kivik/v4"

	"github.com/labstack/echo/v4"

	"app/internal/ical"
)

const sourceURL = "https://raw.githubusercontent.com/bigfoott/ScrapedDuck/data/events.min.json"
const mimeICS = "text/calendar"

func (s *Service) CalendarGet(c echo.Context) error {
	ctx := c.Request().Context()

	var id string
	if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
		return err
	}

	var cfg *CalendarConfig
	if err := s.DB.Get(ctx, id).ScanDoc(&cfg); err != nil {
		if kivik.HTTPStatus(err) == http.StatusNotFound {
			return echo.ErrNotFound
		}

		return fmt.Errorf("failed to get document: %w", err)
	}

	events, err := s.getEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	tz, err := validateTZ(cfg.Timezone)
	if err != nil {
		return err
	}

	options := ical.GenerateICalOptions{
		Now:          time.Now(),
		TZ:           tz,
		IncludeTypes: cfg.IncludeEvents,
		ExcludeTypes: cfg.ExcludeEvents,
	}

	ical, err := ical.GenerateICal(events, options)
	if err != nil {
		return fmt.Errorf("failed to generate ics: %w", err)
	}

	c.Response().Header().Set(echo.HeaderContentType, mimeICS)
	return c.String(http.StatusOK, ical.Serialize())
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
