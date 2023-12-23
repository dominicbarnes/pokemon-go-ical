package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
	_ "time/tzdata"

	ics "github.com/arran4/golang-ical"
	"github.com/segmentio/cli"
	"gopkg.in/yaml.v3"

	"app/internal/ical"
)

const defaultSourceURL = "https://raw.githubusercontent.com/bigfoott/ScrapedDuck/data/events.min.json"

type config struct {
	Input  string `flag:"-i,--input" help:"The input configuration file" default:"config.yaml"`
	Cache  string `flag:"-c,--cache" help:"The cache location" default:"cache.yaml"`
	Output string `flag:"-o,--output" help:"The output directory" default:"./data"`
}

type calendar struct {
	Timezone string   `yaml:"timezone"`
	Include  []string `yaml:"include"`
	Exclude  []string `yaml:"exclude"`
}

type cache struct {
	Config string `yaml:"config"`
	Events string `yaml:"events"`
}

func main() {
	cli.Exec(cli.Command(func(ctx context.Context, config config) error {
		start := time.Now()
		slog.Info("started generator", slog.Time("started_at", start))

		if err := os.MkdirAll(config.Output, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		previous, err := readCache(ctx, config.Cache)
		if err != nil {
			return fmt.Errorf("failed to get cache: %w", err)
		}

		configHash, calendars, err := getCalendars(ctx, config.Input)
		if err != nil {
			return fmt.Errorf("failed to get calendars: %w", err)
		}
		slog.Info("found calendars", slog.Int("calendar_count", len(calendars)))

		eventsHash, events, err := getEvents(ctx, config.Input)
		if err != nil {
			return fmt.Errorf("failed to get events: %w", err)
		}
		slog.Info("found events", slog.Int("event_count", len(events)))

		if previous.Config == configHash && previous.Events == eventsHash {
			slog.Info("no changes detected in calendar config or events")
			return nil
		}

		for name, calendar := range calendars {
			logger := slog.With(slog.String("name", name))
			logger.Debug("generating calendar")

			tz, err := getTimezone(ctx, calendar.Timezone)
			if err != nil {
				return fmt.Errorf("failed to get timezone: %w", err)
			}

			options := ical.GenerateICalOptions{
				Now:          start,
				TZ:           tz,
				IncludeTypes: calendar.Include,
				ExcludeTypes: calendar.Exclude,
			}

			ics, err := ical.GenerateICal(events, options)
			if err != nil {
				return fmt.Errorf("failed to generate ics: %w", err)
			}
			logger.Info("generated calendar",
				slog.String("timezone", tz.String()),
				slog.Any("include", calendar.Include),
				slog.Any("exclude", calendar.Exclude),
				slog.Int("events", len(ics.Events())),
			)

			output := filepath.Join(config.Output, name+".ics")
			if err := writeCalendar(ctx, ics, output); err != nil {
				return fmt.Errorf("failed to write %s: %w", output, err)
			}
			logger.Info("output calendar", slog.String("path", output))
		}

		if err := writeCache(ctx, config.Cache, configHash, eventsHash); err != nil {
			return fmt.Errorf("failed to write cache: %w", err)
		}

		slog.Info("done")
		return nil
	}))
}

func readCache(ctx context.Context, path string) (*cache, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return new(cache), nil
		}

		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var output cache
	if err := yaml.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("failed to decode cache: %w", err)
	}

	return &output, nil
}

func getEvents(ctx context.Context, input string) (string, []ical.LeekDuckEvent, error) {
	res, err := http.Get(defaultSourceURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to download events: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read events: %w", err)
	}

	events, err := decodeEvents(ctx, bytes.NewBuffer(data))
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode events: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	return hash, events, nil
}

func decodeEvents(ctx context.Context, data io.Reader) ([]ical.LeekDuckEvent, error) {
	var ee []ical.LeekDuckEvent
	d := json.NewDecoder(data)
	if err := d.Decode(&ee); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}
	return ee, nil
}

func getTimezone(ctx context.Context, input string) (*time.Location, error) {
	tz, err := time.LoadLocation(input)
	if err != nil {
		return nil, fmt.Errorf("invalid location: %w", err)
	}
	return tz, nil
}

func getCalendars(ctx context.Context, config string) (string, map[string]calendar, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var output map[string]calendar
	if err := yaml.Unmarshal(data, &output); err != nil {
		return "", nil, fmt.Errorf("failed to decode config as YAML: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	return hash, output, nil
}

func writeCalendar(ctx context.Context, cal *ics.Calendar, output string) error {
	ics := cal.Serialize()
	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	if _, err := fmt.Fprint(f, ics); err != nil {
		return fmt.Errorf("could not write calendar: %w", err)
	}
	return nil
}

func writeCache(ctx context.Context, path, config, events string) error {
	output := cache{
		Config: config,
		Events: events,
	}
	data, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("failed to encode cache: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
