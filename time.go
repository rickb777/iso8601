package iso8601

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	RFC3339      = time.RFC3339
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"
	RFC3339Nano  = time.RFC3339Nano
)

// MarshalTextFormat is the rendering format used by MarshalText and MarshalJSON.
// The default, RFC3339Nano, is suitable in most cases. However, if reduced
// precision is required (e.g. when communicating with legacy systems such as
// Salesforce), then this should be set to RFC3339Micro, RFC3339Milli or RFC3339.
//
// This must not be altered concurrently.
//
// If there is a need to marshal using various precision formats, not just one,
// then the values should be rounded using Truncate. Round can also be used, but
// note that the rounding might allow more digits to be sent.
var MarshalTextFormat = RFC3339Nano

var _ json.Unmarshaler = &Time{}

// Date returns the Time corresponding to
//
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
//
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice. In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// Date returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// Date panics if loc is nil.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Of(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

// Of is a construction helper that wraps time.Time as a Time.
func Of(t time.Time) Time {
	return Time{Time: t}
}

// Time adapts time.Time for formatting and parsing ISO-8061 dates,
// especially as a JSON string.
type Time struct {
	time.Time
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Truncate(d time.Duration) Time {
	return Of(t.Time.Truncate(d))
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Round(d time.Duration) Time {
	return Of(t.Time.Round(d))
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in ISO-8601 / RFC 3339 format, with sub-second
// precision controlled by MarshalTextFormat.
func (t Time) MarshalText() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(MarshalTextFormat))
	return t.AppendFormat(b, MarshalTextFormat), nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in ISO-8601 / RFC 3339 format, with sub-second
// precision controlled by MarshalTextFormat.
func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(MarshalTextFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, MarshalTextFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *Time) UnmarshalText(data []byte) error {
	// Fractional seconds are handled implicitly by Parse.
	tt, err := Parse(data)
	if err == nil {
		*t = Of(tt)
	}
	return err
}

// UnmarshalJSON decodes a JSON string or null into a iso8601 time
func (t *Time) UnmarshalJSON(b []byte) error {
	// Do not process null types
	if null(b) {
		return nil
	}
	if len(b) > 0 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	} else {
		return ErrNotString
	}
	var err error
	t.Time, err = Parse(b)
	return err
}

// null returns true if the given byte slice is a JSON null.
// This is about 3x faster than `bytes.Compare`.
func null(b []byte) bool {
	if len(b) != 4 {
		return false
	}
	if b[0] != 'n' && b[1] != 'u' && b[2] != 'l' && b[3] != 'l' {
		return false
	}
	return true
}

// String renders the time in ISO-8601 format (using RFC3339Nano).
func (t Time) String() string {
	// time.RFC3339Nano is one of several permitted ISO-8601 formats.
	return t.Format(RFC3339Nano)
}

// Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec
// nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range
// [0, 999999999]. Not all sec values have a corresponding time value. One such value
// is 1<<63-1 (the largest int64 value).
//func Unix(sec int64, nsec int64) Time {
//	return Of(time.Unix(sec, nsec))
//}

// UnixMicro returns the local Time corresponding to the given Unix time, usec
// microseconds since January 1, 1970 UTC.
//func UnixMicro(usec int64) Time {
//	return Of(time.UnixMicro(usec))
//}

// UnixMilli returns the local Time corresponding to the given Unix time, msec milliseconds
// since January 1, 1970 UTC.
//func UnixMilli(msec int64) Time {
//	return Of(time.UnixMilli(msec))
//}

// Add returns the time t+d.
//func (t Time) Add(d time.Duration) Time {
//	return Of(t.Time.Add(d))
//}

// AddDate returns the time corresponding to adding the given number of years, months,
// and days to t. For example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.
//
// Note that dates are fundamentally coupled to timezones, and calendrical periods like days don't
// have fixed durations. AddDate uses the Location of the Time value to determine these durations.
// That means that the same AddDate arguments can produce a different shift in absolute time
// depending on the base Time value and its Location. For example, AddDate(0, 0, 1) applied to
// is a 24 hour shift. In others it's a 23 hour shift due to daylight savings time transitions.
//
// AddDate normalizes its result in the same way that Date does, so, for example, adding one
// month to October 31 yields December 1, the normalized form for November 31.
//func (t Time) AddDate(years int, months int, days int) Time {
//	return Of(t.Time.AddDate(years, months, days))
//}

// In returns a copy of t representing the same time instant, but with the copy's location
// information set to loc for display purposes.
//func (t Time) In(loc *time.Location) Time {
//	return Of(t.Time.In(loc))
//}

// Local returns t with the location set to local time.
//func (t Time) Local() Time {
//	return Of(t.Time.Local())
//}

// UTC returns t with the location set to UTC.
//func (t Time) UTC() Time {
//	return Of(t.Time.UTC())
//}

// ZoneBounds returns the bounds of the time zone in effect at time t. The zone begins at
// start and the next zone begins at end. If the zone begins at the beginning of time,
// start will be returned as a zero Time. If the zone goes on forever, end will be returned
// as a zero Time. The Location of the returned times will be the same as t.
//func (t Time) ZoneBounds() (start, end Time) {
//	s, e := t.Time.ZoneBounds()
//	return Of(s), Of(e)
//}
