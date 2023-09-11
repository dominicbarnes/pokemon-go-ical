package main

import (
	"testing"
	"time"

	"github.com/dominicbarnes/got"
	"github.com/stretchr/testify/require"
)

var epoch = time.Date(2023, time.September, 11, 6, 11, 0, 0, time.UTC)

func TestGenerateICal(t *testing.T) {
	type Test struct {
		Events  []LeekDuckEvent     `testdata:"events.json"`
		Options GenerateICalOptions `testdata:"options.json"`
	}

	type Expected struct {
		Output string `testdata:"expected/output.ics"`
	}

	suite := got.TestSuite{
		Dir: "testdata/generate-ical",
		TestFunc: func(t *testing.T, c got.TestCase) {
			var test Test
			var expected Expected
			c.Load(t, &test, &expected)

			test.Options.Now = epoch
			test.Options.TZ = time.Local

			actual, err := GenerateICal(test.Events, test.Options)
			require.NoError(t, err)

			got.Assert(t, c.Dir, &Expected{Output: actual.Serialize()})
		},
	}

	suite.Run(t)
}
