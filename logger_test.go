package logger

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
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
		"\x1b[1mTEST>\x1b[0m: %s: test number 1", false},
	{logFmt, colorPrefix, date, Ldate, "test number 2",
		"\x1b[1m\x1b[31mTEST>\x1b[0m: %s: test number 2", false},
	{logFmt, ">>>", time.Kitchen, Ldate | Lshortfile , "test number 3",
		">>>: %s: logger_test.go:56: test number 3", false},
}

func TestOutput(t *testing.T) {
	for i, k := range outputTests {
		// var buf bytes.Buffer
		var buf bytes.Buffer
		log := New(&buf, DEBUG)
		log.Prefix = k.prefix
		log.DateFormat = k.dateFormat
		log.Flags = k.flags
		d := time.Now().Format(log.DateFormat)
		err := log.Fprint(1, k.text, &buf)
		want := fmt.Sprintf(k.want, d)
		if buf.String() != want || err != nil {
			t.Errorf("Print test %d failed, got \"%q\" want "+
				"\"%q\"", i, buf.String(), want)
		}
		log = nil
	}
}
