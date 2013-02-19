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
)

var outputTests = []struct {
	template   string
	prefix     string
	logPrefix  logPrefix
	dateFormat string
	flags      int
	text       string
	want       string
	wantErr    bool
}{
	{logFmt, boldPrefix, PrintPrefix, date, LstdFlags, "test number 1",
		"\x1b[1mTEST>\x1b[0m: %s: test number 1", false},
	{logFmt, colorPrefix, PrintPrefix, date, LstdFlags, "test number 2",
		"\x1b[1m\x1b[31mTEST>\x1b[0m: %s: test number 2", false},
	{logFmt, ">>>", PrintPrefix, time.Kitchen, Ldate | Lshortfile,
		"test number 3", ">>>: %s: logger_test.go:67: test number 3",
		false},
	{logFmt, AnsiEscape(BOLD, ">>>", OFF), PrintPrefix, date, Ldate,
		"test number 4", ">>>: %s: test number 4", false},
	{logFmt, defColorPrefix, DebugPrefix, time.RubyDate, LstdFlags,
		"test number 5", "\x1b[1m\x1b[32m>>>\x1b[0m: \x1b[1m\x1b[37m[DEBUG]\x1b[0m: %s: test number 5", false},
	{logFmt, defColorPrefix, InfoPrefix, time.RubyDate, LstdFlags,
		"test number 6", "\x1b[1m\x1b[32m>>>\x1b[0m: \x1b[1m\x1b[32m[INFO]\x1b[0m: %s: test number 6", false},
	{logFmt, defColorPrefix, WarningPrefix, time.RubyDate, LstdFlags,
		"test number 7", "\x1b[1m\x1b[32m>>>\x1b[0m: \x1b[1m\x1b[33m[WARNING]\x1b[0m: %s: test number 7", false},
	{logFmt, defColorPrefix, ErrorPrefix, time.RubyDate, LstdFlags,
		"test number 8", "\x1b[1m\x1b[32m>>>\x1b[0m: \x1b[1m\x1b[35m[ERROR]\x1b[0m: %s: test number 8", false},
	{logFmt, defColorPrefix, CriticalPrefix, time.RubyDate, LstdFlags,
		"test number 9", "\x1b[1m\x1b[32m>>>\x1b[0m: \x1b[1m\x1b[31m[CRITICAL]\x1b[0m: %s: test number 9", false},
}

func TestOutput(t *testing.T) {
	for i, k := range outputTests {
		var buf bytes.Buffer
		log := New(&buf, DEBUG)
		log.Prefix = k.prefix
		log.DateFormat = k.dateFormat
		log.Flags = k.flags
		d := time.Now().Format(log.DateFormat)
		err := log.Fprint(k.logPrefix, 1, k.text, &buf)
		want := fmt.Sprintf(k.want, d)
		if buf.String() != want || err != nil && !k.wantErr {
			t.Errorf("Print test %d failed, \ngot:  %q\nwant: "+
				"%q", i, buf.String(), want)
			continue
		}
		fmt.Printf("Test %d OK: %s\n", i, buf.String())
	}
}
