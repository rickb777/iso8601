package iso8601

import (
	"strings"
	"testing"
	"time"

	"github.com/rickb777/expect"
)

type TestCase struct {
	Using string

	Year  int
	Month time.Month
	Day   int

	Hour        int
	Minute      int
	Second      int
	MilliSecond int

	Zone float64
}

func TestParse_ok(t *testing.T) {
	var goodCases = []TestCase{
		{
			Using: "2017-04-24T09:41:34.502+0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        1,
		},
		{
			Using: "2017-04-24T09:41+0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41,
			Zone: 1,
		},
		{
			Using: "2017-04-24T09+0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9,
			Zone: 1,
		},
		{
			Using: "2017-04-24T",
			Year:  2017, Month: 4, Day: 24,
		},
		{
			Using: "2017-04-24",
			Year:  2017, Month: 4, Day: 24,
		},
		{
			Using: "2017-04-24T09:41:34+0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			Zone: 1,
		},
		{
			Using: "2017-04-24T09:41:34.502-0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        -1,
		},
		{
			Using: "2017-04-24T09:41:34.502-01:00",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        -1,
		},
		{
			Using: "2017-04-24T09:41-01:00",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41,
			Zone: -1,
		},
		{
			Using: "2017-04-24T09-01:00",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9,
			Zone: -1,
		},
		{
			Using: "2017-04-24T09:41:34-0100",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			Zone: -1,
		},
		{
			Using: "2017-04-24T09:41:34.502Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			Zone: 0,
		},
		{
			Using: "2017-04-24T09:41Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41,
			Zone: 0,
		},
		{
			Using: "2017-04-24T09Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9,
			Zone: 0,
		},
		{
			Using: "2017-04-24T09:41:34.089",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 89,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41,
			Zone: 0,
		},
		{
			Using: "2017-04-24T09",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9,
			Zone: 0,
		},
		{
			Using: "2017-04-24T09:41:34.009",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 9,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34.893",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 893,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34.89312523Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 893,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34.502-0530",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        -5.5,
		},
		{
			Using: "2017-04-24T09:41:34.502+0530",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        5.5,
		},
		{
			Using: "2017-04-24T09:41:34.502+05:30",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        5.5,
		},

		{
			Using: "2017-04-24T09:41:34.502+05:45",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        5.75,
		},
		{
			Using: "2017-04-24T09:41:34.502+00",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34.502+0000",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        0,
		},
		{
			Using: "2017-04-24T09:41:34.502+00:00",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        0,
		},
	}

	for _, c := range goodCases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := ParseString(c.Using) // also encompasses testing of Parse([]byte).
			if err != nil {
				t.Fatal(err)
			}

			expect.Number(d.Year()).ToBe(t, c.Year)
			expect.Number(d.Month()).ToBe(t, c.Month)
			expect.Number(d.Day()).ToBe(t, c.Day)
			expect.Number(d.Hour()).ToBe(t, c.Hour)
			expect.Number(d.Minute()).ToBe(t, c.Minute)
			expect.Number(d.Second()).ToBe(t, c.Second)
			expect.Number(d.Nanosecond()/1000000).ToBe(t, c.MilliSecond)

			_, z := d.Zone()
			expect.Number(float64(z)/3600).ToBe(t, c.Zone)
		})
	}
}

type ErrorCase struct {
	Using   string
	Message string
}

func TestParse_error(t *testing.T) {
	var errorCases = []ErrorCase{
		// Invalid Parse Test Cases
		{
			Using:   "2017-04-24T09:41:34.502-00",
			Message: `Cannot parse "-00": invalid zone`,
		},
		{
			Using:   "2017-04-24T09:41:34.502-0000",
			Message: `Cannot parse "-0000": invalid zone`,
		},
		{
			Using:   "2017-04-24T09:41:34.502-00:00",
			Message: `Cannot parse "-00:00": invalid zone`,
		},

		// Invalid Range Test Cases
		{
			Using:   "2017-00-01T00:00:00.000+00:00",
			Message: `Cannot parse "2017-00-01T00:00:00.000+00:00": month 0 is not in range 1-12`,
		},
		{
			Using:   "2017-13-01T00:00:00.000+00:00",
			Message: `Cannot parse "2017-13-01T00:00:00.000+00:00": month 13 is not in range 1-12`,
		},

		{
			Using:   "2017-01-00T00:00:00.000+00:00",
			Message: `Cannot parse "2017-01-00T00:00:00.000+00:00": day 0 is not in range 1-31`,
		},
		{
			Using:   "2017-01-32T00:00:00.000+00:00",
			Message: `Cannot parse "2017-01-32T00:00:00.000+00:00": day 32 is not in range 1-31`,
		},
		{
			Using:   "2019-02-29T00:00:00.000+00:00",
			Message: `Cannot parse "2019-02-29T00:00:00.000+00:00": day 29 is not in range 1-28`,
		},
		{
			Using:   "2020-02-30T00:00:00.000+00:00", // Leap year
			Message: `Cannot parse "2020-02-30T00:00:00.000+00:00": day 30 is not in range 1-29`,
		},

		{
			Using:   "2017-01-01T24:00:00.000+00:00",
			Message: `Cannot parse "2017-01-01T24:00:00.000+00:00": hour 24 is not in range 0-23`,
		},

		{
			Using:   "2017-01-01T00:60:00.000+00:00",
			Message: `Cannot parse "2017-01-01T00:60:00.000+00:00": minute 60 is not in range 0-59`,
		},

		{
			Using:   "2017-01-01T00:00:60.000+00:00",
			Message: `Cannot parse "2017-01-01T00:00:60.000+00:00": second 60 is not in range 0-59`,
		},
	}

	for _, c := range errorCases {
		t.Run(c.Using, func(t *testing.T) {
			_, err := ParseString(c.Using) // also encompasses testing of Parse([]byte).
			if err == nil {
				t.Fatalf("Expected error containing %q", c.Message)
			} else {
				if !strings.Contains(err.Error(), c.Message) {
					t.Errorf("Expected error message %q to contain %q", err.Error(), c.Message)
				}
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	x := []byte("2017-04-24T09:41:34.502Z")
	for i := 0; i < b.N; i++ {
		_, err := Parse(x)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestParseISOZone(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expect.Any(ParseISOZone([]byte("Z"))).ToBe(t, time.UTC)
		expect.Any(ParseISOZone([]byte("+00:00"))).ToBe(t, time.FixedZone("+00:00", 0))
		expect.Any(ParseISOZone([]byte("+05:00"))).ToBe(t, time.FixedZone("+05:00", 5*3600)) // New York
		expect.Any(ParseISOZone([]byte("+05:01:01"))).ToBe(t, time.FixedZone("+05:01:01", 5*3600+61))
		expect.Any(ParseISOZone([]byte("-01:00"))).ToBe(t, time.FixedZone("-01:00", -3600))
		expect.Any(ParseISOZone([]byte("\u221201:00"))).ToBe(t, time.FixedZone("\u221201:00", -3600))
		expect.Any(ParseISOZone([]byte("+0100"))).ToBe(t, time.FixedZone("+0100", 3600))
		expect.Any(ParseISOZone([]byte("-0100"))).ToBe(t, time.FixedZone("-0100", -3600))
		expect.Any(ParseISOZone([]byte("+01"))).ToBe(t, time.FixedZone("+01", 3600))
		expect.Any(ParseISOZone([]byte("-01"))).ToBe(t, time.FixedZone("-01", -3600))
		expect.Any(ParseISOZone([]byte("+03"))).ToBe(t, time.FixedZone("+03", 3*3600))
		expect.Any(ParseISOZone([]byte("-03"))).ToBe(t, time.FixedZone("-03", -3*3600))
		expect.Any(ParseISOZone([]byte("+0030"))).ToBe(t, time.FixedZone("+0030", 1800))
		expect.Any(ParseISOZone([]byte("-0030"))).ToBe(t, time.FixedZone("-0030", -1800))
	})

	t.Run("error", func(t *testing.T) {
		expect.Error(ParseISOZone([]byte("+0"))).ToContain(t, "iso8601: Zone information is too short")
		expect.Error(ParseISOZone([]byte("1"))).ToContain(t, `iso8601: Cannot parse "1": invalid zone at '1'`)
		expect.Error(ParseISOZone([]byte("+12345678"))).ToContain(t, "iso8601: Zone information is too long")
		expect.Error(ParseISOZone([]byte("0100"))).ToContain(t, `iso8601: Cannot parse "0100": invalid zone at '0'`)
		expect.Error(ParseISOZone([]byte("-0000"))).ToContain(t, "iso8601: Cannot parse \"-0000\": invalid zone")
		expect.Error(ParseISOZone([]byte("-0:10"))).ToContain(t, `iso8601: Cannot parse "-0:10": invalid zone at ':'`)
		expect.Error(ParseISOZone([]byte("-01:0"))).ToContain(t, "iso8601: Cannot parse \"-01:0\": invalid zone")
		expect.Error(ParseISOZone([]byte("-foo"))).ToContain(t, `iso8601: Cannot parse "-foo": invalid zone at 'f'`)
		expect.Error(ParseISOZone([]byte{0xAA, 0xBB})).ToContain(t, `iso8601: Cannot parse "\xaa\xbb": invalid zone at '?'`)
	})
}
