package iso8601

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/rickb777/expect"
)

type TestAPIResponse struct {
	Ptr  *Time
	Nptr Time
}

type TestStdLibAPIResponse struct {
	Ptr  *time.Time
	Nptr time.Time
}

var ShortTest = TestCase{
	Using: "2001-11-13",
	Year:  2001, Month: 11, Day: 13,
}

var StructTestData = []byte(`
{
  "Ptr": "2017-04-26T11:13:04+01:00",
  "Nptr": "2017-04-26T11:13:04+01:00"
}
`)

var NullTestData = []byte(`
{
  "Ptr": null,
  "Nptr": null
}
`)

var ZeroedTestData = []byte(`
{
  "Ptr": "0001-01-01",
  "Nptr": "0001-01-01"
}
`)

var StructTest = TestCase{
	Year: 2017, Month: 04, Day: 26,
	Hour: 11, Minute: 13, Second: 04,
	Zone: 1,
}

func TestTime_Unmarshaling(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		var b = []byte(`"2001-11-13"`)

		tn := new(Time)
		if err := tn.UnmarshalJSON(b); err != nil {
			t.Fatal(err)
		}

		expect.Number(tn.Year()).ToBe(t, ShortTest.Year)
		expect.Number(tn.Month()).ToBe(t, ShortTest.Month)
		expect.Number(tn.Day()).ToBe(t, ShortTest.Day)

		expect.Any(tn.UnmarshalJSON([]byte(`2001-11-13`))).ToBe(t, ErrNotString)
	})

	t.Run("struct", func(t *testing.T) {
		resp := new(TestAPIResponse)
		expect.Error(json.Unmarshal(StructTestData, resp)).Not().ToHaveOccurred(t)

		stdlibResp := new(TestStdLibAPIResponse)
		expect.Error(json.Unmarshal(StructTestData, stdlibResp)).Not().ToHaveOccurred(t)

		t.Run("stblib parity", func(t *testing.T) {
			if !resp.Ptr.Time.Equal(*stdlibResp.Ptr) || !resp.Nptr.Time.Equal(stdlibResp.Nptr) {
				t.Fatalf("Parsed time values are not equal to standard library implementation")
			}
		})

		t.Run("ptr", func(t *testing.T) {
			expect.Number(resp.Ptr.Year()).ToBe(t, StructTest.Year)
			expect.Number(resp.Ptr.Day()).ToBe(t, StructTest.Day)
			expect.Number(resp.Ptr.Second()).ToBe(t, StructTest.Second)
		})

		t.Run("noptr", func(t *testing.T) {
			expect.Number(resp.Nptr.Year()).ToBe(t, StructTest.Year)
			expect.Number(resp.Nptr.Day()).ToBe(t, StructTest.Day)
			expect.Number(resp.Nptr.Second()).ToBe(t, StructTest.Second)
		})
	})

	t.Run("null", func(t *testing.T) {
		resp := new(TestAPIResponse)
		expect.Error(json.Unmarshal(NullTestData, resp)).Not().ToHaveOccurred(t)
	})

	t.Run("time zeroed", func(t *testing.T) {
		resp := new(TestAPIResponse)
		expect.Error(json.Unmarshal(ZeroedTestData, resp)).Not().ToHaveOccurred(t)
	})

	t.Run("reparse", func(t *testing.T) {
		s := time.Now().UTC()
		data := []byte(s.Format(time.RFC3339Nano))
		n, err := Parse(data)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Any(s).ToBe(t, n.Time)
	})

	t.Run("string", func(t *testing.T) {
		t1 := time.Now().UTC()
		s := Time{Time: t1}.String()
		expect.String(s).ToBe(t, t1.Format(time.RFC3339Nano))
	})
}

func TestTime_Marshaling(t *testing.T) {
	t9 := Date(2017, 4, 26, 11, 13, 4, 123456789, time.UTC)

	cases := []struct {
		format     string
		resolution time.Duration
		expected   string
	}{
		{
			format:     RFC3339,
			resolution: time.Second,
			expected:   "2017-04-26T11:13:04Z",
		},
		{
			format:     RFC3339Milli,
			resolution: time.Millisecond,
			expected:   "2017-04-26T11:13:04.123Z",
		},
		{
			format:     RFC3339Micro,
			resolution: time.Microsecond,
			expected:   "2017-04-26T11:13:04.123456Z",
		},
		{
			format:     RFC3339Nano,
			resolution: time.Nanosecond,
			expected:   "2017-04-26T11:13:04.123456789Z",
		},
	}

	t.Run("text marshal/unmarshal", func(t *testing.T) {
		for _, c := range cases {
			MarshalTextFormat = c.format

			b, err := xml.Marshal(t9)
			expect.String(b, err).ToEqual(t, fmt.Sprintf("<Time>%s</Time>", c.expected))

			tn := Time{}
			err = xml.Unmarshal(b, &tn)

			expect.Any(tn, err).ToBe(t, t9.Truncate(c.resolution))
		}
	})

	t.Run("JSON marshal/unmarshal", func(t *testing.T) {
		for _, c := range cases {
			MarshalTextFormat = c.format

			b, err := json.Marshal(t9)
			expect.String(b, err).ToEqual(t, fmt.Sprintf("%q", c.expected))

			var tn Time
			err = tn.UnmarshalJSON(b)

			expect.Any(tn, err).ToBe(t, t9.Truncate(c.resolution))
		}
	})
}

func TestTime_Decorators(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	expect.Error(err).ToBeNil(t)

	t9 := Date(2017, 4, 26, 11, 13, 4, 123456789, ny)

	t.Run("Unix", func(t *testing.T) {
		r := Unix(101, 987)
		expect.Any(r.Time).ToBe(t, time.Unix(101, 987))
	})

	t.Run("UnixMicro", func(t *testing.T) {
		r := UnixMicro(123456789)
		expect.Any(r.Time).ToBe(t, time.UnixMicro(123456789))
	})

	t.Run("UnixMilli", func(t *testing.T) {
		r := UnixMilli(123456789)
		expect.Any(r.Time).ToBe(t, time.UnixMilli(123456789))
	})

	t.Run("Local", func(t *testing.T) {
		r := t9.Local()
		expect.Any(r.Time).ToBe(t, t9.Time.Local())
	})

	t.Run("UTC", func(t *testing.T) {
		r := t9.UTC()
		expect.Any(r.Time).ToBe(t, t9.Time.UTC())
	})

	t.Run("In", func(t *testing.T) {
		r := t9.In(time.Local)
		expect.Any(r.Time).ToBe(t, t9.Time.In(time.Local))
	})

	t.Run("ZoneBounds", func(t *testing.T) {
		s, e := t9.ZoneBounds()
		expStart, expEnd := t9.Time.ZoneBounds()
		expect.Any(s.Time).ToBe(t, expStart)
		expect.Any(e.Time).ToBe(t, expEnd)
	})

	t.Run("Truncate", func(t *testing.T) {
		tr := t9.Truncate(time.Microsecond)
		expect.Any(tr.Time).ToEqual(t, t9.Time.Truncate(time.Microsecond))
	})

	t.Run("Round", func(t *testing.T) {
		r := t9.Round(time.Microsecond)
		expect.Any(r.Time).ToBe(t, t9.Time.Round(time.Microsecond))
	})

	t.Run("Add", func(t *testing.T) {
		r := t9.Add(123 * time.Microsecond)
		expect.Any(r.Time).ToBe(t, t9.Time.Add(123*time.Microsecond))
	})

	t.Run("AddDate", func(t *testing.T) {
		r := t9.AddDate(1, 2, 3)
		expect.Any(r.Time).ToBe(t, t9.Time.AddDate(1, 2, 3))
	})

}

func BenchmarkCheckNull(b *testing.B) {
	var n = []byte("null")

	b.Run("compare", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bytes.Compare(n, n)
		}
	})
	b.Run("exact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			null(n)
		}
	})
}
