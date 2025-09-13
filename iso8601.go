// Package iso8601 is a utility for parsing ISO8601 datetime strings into native Go times.
// The standard library's RFC3339 reference layout can be too strict for working with 3rd party APIs,
// especially ones written in other languages.
//
// Use the provided `Time` structure instead of the default `time.Time` to provide ISO8601 support for JSON responses.
package iso8601

import (
	"time"
	"unicode/utf8"
)

const (
	year uint = iota
	month
	day
	hour
	minute
	second
	millisecond
)

const (
	// charStart is the binary position of the character `0`
	charStart int = '0'
)

// ParseISOZone parses the zone information in an ISO8061 date string.
// Most timezones use only hours and minutes; seconds are also supported but
// not fractions of seconds.
//
// The input is expected to match one of:
//
//	 Z
//		+hh
//		-hh
//		+hh:mm
//		-hh:mm
//		+hh:mm:ss
//		-hh:mm:ss
//
// The leading character can be plus + (u002B), hyphen - (u002D) or minus − (u2212).
// Examples: Z, +03, -0100, +02:00, +01:45:30
func ParseISOZone(inp []byte) (*time.Location, error) {
	var neg bool

	r, i := utf8.DecodeRune(inp)
	switch r {
	case 'Z':
		return time.UTC, nil
	case '+':
	case '-', '\u2212':
		neg = true
	default:
		if r == utf8.RuneError {
			return nil, newUnexpectedCharacterError('?')
		}
		return nil, newUnexpectedCharacterError(r)
	}

	if len(inp) < 3 {
		return nil, ErrZoneTooShort
	}

	var offset int
	number := inp[i:]

	var z, digits int
	var multiplier = 3600 // start with initial multiplier of hours
	for i = 0; i < len(number); i++ {
		if digits > 2 {
			return nil, ErrZoneTooLong
		} else if i == 2 || i == 5 { // next multiplier
			offset += z * multiplier
			multiplier /= 60 // multiplier for minutes or seconds
			z = 0
			digits = 0
		} else { // next digit
			z = z * 10
		}

		switch number[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			z += int(number[i]) - charStart
			digits++
		case ':':
			if i != 2 && i != 5 {
				return nil, newUnexpectedCharacterError(rune(number[i]))
			}
			digits = 0
		default:
			return nil, newUnexpectedCharacterError(rune(number[i]))
		}
	}

	offset += z * multiplier

	if digits != 2 {
		return nil, ErrInvalidZone
	}

	if neg {
		offset = -offset
	}

	if neg && offset == 0 {
		return nil, ErrInvalidZone
	}

	return time.FixedZone(string(inp), offset), nil
}

// Parse parses an ISO8601 compliant date-time byte slice into a time.Time object.
// If any component of an input date-time is not within the expected range then an *iso8601.RangeError is returned.
func Parse(inp []byte) (Time, error) {
	var (
		Y         int
		M         int
		d         int
		h         int
		m         int
		s         int
		fraction  int
		nfraction = 1 // counts amount of precision for the second fraction
	)

	// Always assume UTC by default
	var loc = time.UTC

	var c int
	var p = year

	var i int

parse:
	for ; i < len(inp); i++ {
		switch inp[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			c = c * 10
			c += int(inp[i]) - charStart

			if p == millisecond {
				nfraction++
			}
		case '-':
			if p < hour {
				switch p {
				case year:
					Y = c
				case month:
					M = c
				default:
					return Time{}, newUnexpectedCharacterError(rune(inp[i]))
				}
				p++
				c = 0
				continue
			}
			fallthrough
		case '+':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				s = c
			case millisecond:
				fraction = c
			default:
				return Time{}, newUnexpectedCharacterError(rune(inp[i]))
			}
			c = 0
			var err error
			loc, err = ParseISOZone(inp[i:])
			if err != nil {
				return Time{}, err
			}
			break parse
		case 'T':
			if p != day {
				return Time{}, newUnexpectedCharacterError(rune(inp[i]))
			}
			d = c
			c = 0
			p++
		case ':':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				m = c
			default:
				return Time{}, newUnexpectedCharacterError(rune(inp[i]))
			}
			c = 0
			p++
		case '.':
			if p != second {
				return Time{}, newUnexpectedCharacterError(rune(inp[i]))
			}
			s = c
			c = 0
			p++
		case 'Z':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				s = c
			case millisecond:
				fraction = int(c)
			default:
				return Time{}, newUnexpectedCharacterError(rune(inp[i]))
			}
			c = 0
			if len(inp) != i+1 {
				return Time{}, ErrRemainingData
			}
		default:
			return Time{}, newUnexpectedCharacterError(rune(inp[i]))
		}
	}

	// Capture remaining data
	// Sometimes a date can end without a non-integer character
	if c > 0 {
		switch p {
		case day:
			d = c
		case hour:
			h = c
		case minute:
			m = c
		case second:
			s = c
		case millisecond:
			fraction = c
		}
	}

	// Get the seconds fraction as nanoseconds
	if fraction < 0 || 1e9 <= fraction {
		return Time{}, ErrPrecision
	}
	scale := 10 - nfraction
	for i := 0; i < scale; i++ {
		fraction *= 10
	}

	switch {
	case M < 1 || M > 12: // Month 1-12
		return Time{}, &RangeError{
			Value:   string(inp),
			Element: "month",
			Given:   M,
			Min:     1,
			Max:     12,
		}
	case d < 1 || d > daysIn(time.Month(M), Y): // Day 1-daysIn(month, year)
		return Time{}, &RangeError{
			Value:   string(inp),
			Element: "day",
			Given:   d,
			Min:     1,
			Max:     daysIn(time.Month(M), Y),
		}
	case h > 23: // Hour 0-23
		return Time{}, &RangeError{
			Value:   string(inp),
			Element: "hour",
			Given:   h,
			Min:     0,
			Max:     23,
		}
	case m > 59: // Minute 0-59
		return Time{}, &RangeError{
			Value:   string(inp),
			Element: "minute",
			Given:   m,
			Min:     0,
			Max:     59,
		}
	case s > 59: // Second 0-59
		return Time{}, &RangeError{
			Value:   string(inp),
			Element: "second",
			Given:   s,
			Min:     0,
			Max:     59,
		}
	}

	return Date(Y, time.Month(M), d, h, m, s, fraction, loc), nil
}

// ParseString parses an ISO8601 compliant date-time string into a time.Time object.
func ParseString(inp string) (Time, error) {
	return Parse([]byte(inp))
}

// String renders the time in ISO-8601 format (using RFC3339Nano).
func (t Time) String() string {
	// time.RFC3339Nano is one of several permitted ISO-8601 formats.
	return t.Format(RFC3339Nano)
}
