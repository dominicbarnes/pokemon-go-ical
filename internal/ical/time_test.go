package ical

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseLeekDuckTime(t *testing.T) {
	spec := []struct {
		input    string
		tz       *time.Location
		expected time.Time
	}{
		{
			input:    "2023-09-01T10:00:00.000",
			tz:       time.Local,
			expected: time.Date(2023, time.September, 1, 10, 00, 00, 000, time.Local),
		},
		{
			input:    "2023-09-08T13:00:00.000-0700",
			tz:       time.Local,
			expected: time.Date(2023, time.September, 8, 13, 00, 00, 000, time.FixedZone("", -7*60*60)),
		},
		{
			input:    "2023-09-22T20:00:00.000Z",
			tz:       time.Local,
			expected: time.Date(2023, time.September, 22, 20, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range spec {
		actual, err := parseLeekDuckEventTime(test.input, test.tz)
		require.NoError(t, err)
		require.Equal(t, test.expected.Format(time.RFC3339Nano), (*actual).Format(time.RFC3339Nano))
	}
}
