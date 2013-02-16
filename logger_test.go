package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestStream(t *testing.T) {
	if out := log.Stream; out != os.Stderr {
		t.Errorf("log.Stream is not stderr by default")
	}
	var buf bytes.Buffer
	log.Stream = &buf
	if out := log.Stream; out != &buf {
		t.Errorf("log.Stream = %p, want %p", out, &buf)
	}
}

var (
	boldPrefix  = AnsiEscape(BOLD, "TEST>", OFF)
	colorPrefix = AnsiEscape(BOLD, RED, "TEST>", OFF)
	date        = "Mon Jan 02 15:04 2006"
	badDate     = "on Jan _2 1504:05 006"
)

var outputTests = []struct {
	template   string
	prefix     string
	dateFormat string
	flags      int
	text       string
	want       string
	wantErr    bool
}{
	{logFmt, boldPrefix, date, LstdFlags, "test number 1",
		"\x1b[1mTEST>\x1b[0m: Mon Jan 02 15:04 2006: test number 1",
		false},
	{logFmt, colorPrefix, date, Ldate, "test number 2",
		"\x1b[1m\x1b[31mTEST>\x1b[0m: Mon Jan 02 15:04 2006: test number 2",
		false},
}

func TestOutput(t *testing.T) {
	for i, k := range outputTests {
		var buf bytes.Buffer
		log := New(&buf, DEBUG)
		log.Prefix = k.prefix
		log.DateFormat = k.dateFormat
		log.Flags = k.flags
		_, err := log.Fprint(1, k.text, &buf)
		if buf.String() != k.want || err != nil {
			t.Errorf("Print test %d failed, got \"%s\" want "+
				"\"%s\"", i, buf.String(), k.want)
		}
		log = nil
	}
}
