package iso8601

import (
	"errors"
	"fmt"
)

var (
	// ErrZoneTooShort indicates too few characters were passed to ParseISOZone.
	ErrZoneTooShort = errors.New("iso8601: Zone information is too short")

	// ErrZoneTooLong indicates too many characters were passed to ParseISOZone.
	ErrZoneTooLong = errors.New("iso8601: Zone information is too long")

	// ErrRemainingData indicates that there is extra data after a `Z` character.
	ErrRemainingData = errors.New("iso8601: Unexpected remaining data after `Z`")

	// ErrNotString indicates that a non string type was passed to the UnmarshalJSON method of `Time`.
	ErrNotString = errors.New("iso8601: Invalid json type (expected string)")

	// ErrPrecision indicates that there was too much precision (characters) given to parse
	// for the fraction of a second of the input time.
	ErrPrecision = errors.New("iso8601: Too many characters in fraction of second precision")
)

func newUnexpectedCharacterError(c rune) error {
	return &UnexpectedCharacterError{Character: c}
}

// UnexpectedCharacterError indicates the parser scanned a character that was not expected at that time.
type UnexpectedCharacterError struct {
	Character rune
}

func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("iso8601: Unexpected character `%c`", e.Character)
}

type SyntaxError struct {
	Value   string
	Element string
	Rune    rune
}

func (e *SyntaxError) Error() string {
	if e.Rune == 0 {
		return fmt.Sprintf("iso8601: Cannot parse %q: invalid %s", e.Value, e.Element)
	}
	return fmt.Sprintf("iso8601: Cannot parse %q: invalid %s at '%c'", e.Value, e.Element, e.Rune)
}

type RangeError struct {
	Value   string
	Element string
	Min     int
	Max     int
	Given   int
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("iso8601: Cannot parse %q: %s %d is not in range %d-%d", e.Value, e.Element, e.Given, e.Min, e.Max)
}
