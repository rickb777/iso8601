A fast ISO8601 date parser for Go

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](http://pkg.go.dev/github.com/rickb777/iso8601)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/iso8601)](https://goreportcard.com/report/github.com/rickb777/iso8601)
[![Build](https://github.com/rickb777/iso8601/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/iso8601/actions)
[![Issues](https://img.shields.io/github/issues/rickb777/iso8601.svg)](https://github.com/rickb777/iso8601/issues)

```
go get github.com/rickb777/iso8601/v3
```

The built-in RFC3333 time layout in Go is too restrictive to support any ISO8601 date-time.

This library parses any ISO8601 date into a native Go time object without regular expressions.

## Usage

```go
import "github.com/rickb777/iso8601/v3"

// iso8601.Time can be used as a drop-in replacement for time.Time with JSON responses
type ExternalAPIResponse struct {
	Timestamp *iso8601.Time
}


func main() {
	// iso8601.ParseString can also be called directly
	t, err := iso8601.ParseString("2020-01-02T16:20:00")
}
```

## Benchmark

```
BenchmarkParse-16        	13364954	        77.7 ns/op	       0 B/op	       0 allocs/op
```

## Release History

  - `3.0.0`

  Parse & ParseString now return iso8601.Time, which emulates time.Time but replaces methods that return Time so that
  it can be a drop-in replacement.
  
  - `2.1.0`

  Added var MarshalTextFormat to allow reduction in precision of XML and JSON values, where this is needed by legacy systems.
  Added methods Round, Truncate, MarshalText, MarshalJSON and UnmarshalText to assist with greater control of precision in marshaled values.
  
  - `2.0.1`

  Fixes the go.mod module.
  
  - `2.0.0` 
  
  Time range validity checking is now equivalent to the standard library. Previous versions would not validate that a given date string was in the expected range. Nor does it support leap seconds (such that the seconds field is `60`), so behaving the same as the [standard library](https://github.com/golang/go/issues/15247)

  Similarly, this version no longer accepts `0000-00-00T00:00:00` as a valid input, even though this can be the zero time representation in other languages.

  - `1.1.0` 
  
  Check for `-0` time zone

  - `1.0.0` 
  
  Initial release
