package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	level := Level()
	if level != WARNING {
		t.Errorf("Level() = %s, want %s", level, WARNING)
	}
	SetLevel(DEBUG)
	level = Level()
	if level != DEBUG {
		t.Errorf("Level() = %s, want %s", level, DEBUG)
	}
}

func TestStream(t *testing.T) {
	if out := Stream(); out != os.Stderr {
		t.Errorf("Output stream is not stderr by default")
	}
	var buf bytes.Buffer
	SetStream(&buf)
	if out := Stream(); out != &buf {
		t.Errorf("SetOutput(&buf) = %p, want %p", out, &buf)
	}
}

var (
	boldPrefix  = AnsiEscape(BOLD, "TEST>", OFF)
	colorPrefix = AnsiEscape(BOLD, RED, "TEST>", OFF)
	date        = "Mon Jan 02 15:04 2006"
	badDate     = "on Jan _2 1504:05 006"
)

var printTests = []struct {
	template   string
	prefix     string
	dateFormat string
	flags      int
	text       string
	want       string
	wantErr    bool
}{
	{logFormat, boldPrefix, date, LstdFlags, "test number 1",
		"\x1b[1mTEST>\x1b[0m: Mon Jan 02 15:04 2006: test number 1",
		false},
	{logFormat, colorPrefix, date, Ldate, "test number 2",
		"\x1b[1m\x1b[31mTEST>\x1b[0m: Mon Jan 02 15:04 2006: test number 2",
		false},
}

func TestOutput(t *testing.T) {
	var buf bytes.Buffer
	for i, k := range printTests {
		log := New(&buf, k.prefix, k.dateFormat, k.template, DEBUG,
			k.flags)
		_, err := log.Output(1, k.text, os.Stdout)
		if buf.String() != k.want || err != nil {
			t.Errorf("Print test %d failed, got \"%s\" want "+
				"\"%s\"", i, buf.String(), k.want)
		}
		log = nil
	}
}
