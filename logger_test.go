// Copyright 2013,2014 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package log

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"text/template"
	"time"

	"github.com/demizer/rgbterm"
)

func TestStream(t *testing.T) {
	var buf bytes.Buffer
	logr := New(LEVEL_CRITICAL, os.Stdout, &buf)
	if out := logr.Streams()[1]; out != &buf {
		t.Errorf("Stream = %p, want %p", out, &buf)
	}
}

func TestMultiStreams(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fPath := filepath.Join(os.TempDir(), fmt.Sprint("go_test_",
		rand.Int()))
	file, err := os.Create(fPath)
	if err != nil {
		t.Error("Create(%q) = %v; want: nil", fPath, err)
	}
	defer file.Close()
	var buf bytes.Buffer
	eLen := 90
	logr := New(LEVEL_DEBUG, file, &buf)
	logr.Debugln("Testing debug output!")
	b := make([]byte, eLen)
	n, err := file.ReadAt(b, 0)
	if n != eLen || err != nil {
		t.Errorf("Read(%d) = %d, %v; want: %d, nil", eLen, n, err,
			eLen)
	}
	if buf.Len() != eLen {
		t.Errorf("buf.Len() = %d; want: %d", buf.Len(), eLen)
	}
}

func TestLongFileFlag(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | LlongFileName)
	Debugln("Test long file flag")
	_, file, _, _ := runtime.Caller(0)
	expect := fmt.Sprintf("[DEBUG] %s: Test long file flag\n", file)
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestShortFileFlag(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | LshortFileName)

	Debugln("Test short file flag")
	_, file, _, _ := runtime.Caller(0)
	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	file = short
	expect := fmt.Sprintf("[DEBUG] %s: Test short file flag\n", file)
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

var date = "Mon 20060102 15:04:05"

var outputTests = []struct {
	template   string
	prefix     string
	level      level
	dateFormat string
	flags      int
	text       string
	expect     string
	expectErr  bool
}{
	// Test with color prefix
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_ALL,
		dateFormat: date,
		flags:      LstdFlags,
		text:       "test number 1",
		// The %s format specifier is the placeholder for the date.
		expect:    "%s \x1b[38;5;46mTEST>\x1b[0;00m test number 1",
		expectErr: false,
	},
	// Test output with coloring turned off
	{
		template:   logFmt,
		prefix:     "TEST>",
		level:      LEVEL_ALL,
		dateFormat: date,
		flags:      Ldate,
		text:       "test number 2",
		expect:     "%s TEST> test number 2",
		expectErr:  false,
	},
	// Test debug output
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_DEBUG,
		dateFormat: time.RubyDate,
		flags:      LstdFlags,
		text:       "test number 3",
		expect:     "%s \x1b[38;5;46mTEST>\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m test number 3",
		expectErr:  false,
	},
	// Test info output
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_INFO,
		dateFormat: time.RubyDate,
		flags:      LstdFlags,
		text:       "test number 4",
		expect:     "%s \x1b[38;5;46mTEST>\x1b[0;00m \x1b[38;5;41m[INFO]\x1b[0;00m test number 4",
		expectErr:  false,
	},
	// Test warning output
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_WARNING,
		dateFormat: time.RubyDate,
		flags:      LstdFlags,
		text:       "test number 5",
		expect:     "%s \x1b[38;5;46mTEST>\x1b[0;00m \x1b[38;5;228m[WARNING]\x1b[0;00m test number 5",
		expectErr:  false,
	},
	// Test error output
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_ERROR,
		dateFormat: time.RubyDate,
		flags:      LstdFlags,
		text:       "test number 6",
		expect:     "%s \x1b[38;5;46mTEST>\x1b[0;00m \x1b[38;5;200m[ERROR]\x1b[0;00m test number 6",
		expectErr:  false,
	},
	// Test critical output
	{
		template:   logFmt,
		prefix:     rgbterm.String("TEST>", 0, 255, 0),
		level:      LEVEL_CRITICAL,
		dateFormat: time.RubyDate,
		flags:      LstdFlags,
		text:       "test number 7",
		expect:     "%s \x1b[38;5;46mTEST>\x1b[0;00m \x1b[38;5;196m[CRITICAL]\x1b[0;00m test number 7",
		expectErr:  false,
	},
	// Test date format
	{
		template:   logFmt,
		prefix:     "::",
		level:      LEVEL_ALL,
		dateFormat: "Mon 20060102 15:04:05",
		flags:      LstdFlags,
		text:       "test number 8",
		expect:     "%s :: test number 8",
		expectErr:  false,
	},

	// FIXME: RE-ADD SUPPORT FOR BOLD!

	// {logFmt, boldPrefix, LEVEL_ALL, date, LstdFlags, "test number 1",
	// "%s \x1b[1mTEST>\x1b[0m test number 1", false},

}

func TestOutput(t *testing.T) {
	for i, k := range outputTests {
		var buf bytes.Buffer
		logr := New(LEVEL_DEBUG, &buf)
		logr.SetPrefix(k.prefix)
		logr.SetDateFormat(k.dateFormat)
		logr.SetFlags(k.flags)
		logr.SetLevel(k.level)
		d := time.Now().Format(logr.DateFormat())
		n, err := logr.Fprint(k.level, 1, k.text, &buf)
		if n != buf.Len() {
			t.Error("Error: ", io.ErrShortWrite)
		}
		expect := fmt.Sprintf(k.expect, d)
		if buf.String() != expect || err != nil && !k.expectErr {
			t.Errorf("Test Number %d\nGot:\t%q\nExpect:\t"+"%q",
				i+1, buf.String(), expect)
			continue
		}
	}
}

func TestLevel(t *testing.T) {
	var buf bytes.Buffer
	logr := New(LEVEL_CRITICAL, &buf)
	logr.Debug("This level should produce no output")
	if buf.Len() != 0 {
		t.Errorf("Debug() produced output at LEVEL_CRITICAL logging level")
	}
	logr.SetLevel(LEVEL_DEBUG)
	logr.Debug("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the LEVEL_DEBUG logging level")
	}
	buf.Reset()
	logr.SetLevel(LEVEL_CRITICAL)
	logr.Println("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the ALL logging level")
	}
	buf.Reset()
	logr.SetLevel(LEVEL_ALL)
	logr.Debug("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the ALL logging level")
	}
}

func TestPrefixNewline(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetFlags(LnoPrefix)
	Print("\n\nThis line should be padded with newlines.\n\n")
	expect := "\n\nThis line should be padded with newlines.\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLdate(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetFlags(LnoPrefix)
	Println("This output should not have a date.")
	expect := "This output should not have a date.\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLfunctionName(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetFlags(LnoPrefix | LfunctionName)
	Println("This output should have a function name.")
	expect := "TestFlagsLfunctionName: This output should have a function name.\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLfunctionNameWithFileName(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetFlags(LnoPrefix | LfunctionName | LshortFileName)
	Print("This output should have a file name and a function name.")
	expect := "logger_test.go: TestFlagsLfunctionNameWithFileName" +
		": This output should have a file name and a function name."
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsNoLansiWithNewlinePadding(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_ALL)
	SetFlags(LnoPrefix)
	Debug("\n\nThis output should be padded with newlines and not colored.\n\n")
	expect := "\n\n[DEBUG] This output should be padded with newlines and not colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLansiWithNewlinePaddingDebug(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_ALL)
	SetFlags(LnoPrefix | Lansi)
	Debug("\n\nThis output should be padded with newlines and colored.\n\n")
	expect := "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m This output should be " +
		"padded with newlines and colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLansiWithNewlinePaddingDebugf(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_ALL)
	SetFlags(LnoPrefix | Lansi)
	Debugf("\n\nThis output should be padded with newlines and %s.\n\n",
		"colored")
	expect := "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m This output should be " +
		"padded with newlines and colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
	buf.Reset()
	Debugf("\n\n##### HELLO %s #####\n\n", "NEWMAN")
	expect = "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m ##### HELLO NEWMAN #####\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestFlagsLansiWithNewlinePaddingDebugln(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_ALL)
	SetFlags(LnoPrefix | Lansi)
	Debugln("\n\nThis output should be padded with newlines and colored.\n\n")
	expect := "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m This output should be " +
		"padded with newlines and colored.\n\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
	buf.Reset()
	Debugln("\n\n", "### HELLO", "NEWMAN", "###", "\n\n")
	expect = "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m  ### HELLO NEWMAN ### \n\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
	buf.Reset()
	Debugln("\n\n### HELLO", "NEWMAN", "###\n\n")
	expect = "\n\n\x1b[38;5;231m[DEBUG]\x1b[0;00m ### HELLO NEWMAN ###\n\n\n"
	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}
}

func TestHeirarchicalPrintln(t *testing.T) {
	var buf bytes.Buffer
	var tBuf bytes.Buffer

	logr := New(LEVEL_ALL, &buf)
	logr.SetFlags(LstdFlags | Lheirarchical)

	now := time.Now()

	logr.Println("\n\nLevel 0 Output 1")
	lvl2 := func() {
		logr.Println("Level 2 Output 1")
		logr.Println("Level 2 Output 2")
	}
	lvl1 := func() {
		logr.Println("Level 1 Output 1")
		logr.Println("Level 1 Output 2")
		lvl2()
	}
	lvl1()

	date = now.Format(std.DateFormat())
	f := struct {
		Date string
	}{Date: date}

	temp := "\n\n{{.Date}} \x1b[38;5;48m::\x1b[0;00m [00] Level 0 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m      [01] Level 1 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m      [01] Level 1 Output 2\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m           [02] Level 2 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m           [02] Level 2 Output 2\n"

	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		t.Fatal(err)
	}
	tmpl.Execute(&tBuf, f)

	if buf.String() != tBuf.String() {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), tBuf.String())
	}
}

func TestHeirarchicalDebugln(t *testing.T) {
	var buf bytes.Buffer
	var tBuf bytes.Buffer

	logr := New(LEVEL_DEBUG, &buf)
	logr.SetFlags(LstdFlags | Lheirarchical)

	now := time.Now()

	logr.Debugln("Level 0 Output 1")
	lvl2 := func() {
		logr.Debugln("Level 2 Output 1")
		logr.Debugln("Level 2 Output 2")
	}
	lvl1 := func() {
		logr.Debugln("Level 1 Output 1")
		logr.Debugln("Level 1 Output 2")
		lvl2()
	}
	lvl1()

	date = now.Format(std.DateFormat())
	f := struct {
		Date string
	}{Date: date}

	temp := "{{.Date}} \x1b[38;5;48m::\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m [00] Level 0 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m      [01] Level 1 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m      [01] Level 1 Output 2\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m           [02] Level 2 Output 1\n" +
		"{{.Date}} \x1b[38;5;48m::\x1b[0;00m \x1b[38;5;231m[DEBUG]\x1b[0;00m           [02] Level 2 Output 2\n"

	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		t.Fatal(err)
	}
	tmpl.Execute(&tBuf, f)

	if buf.String() != tBuf.String() {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), tBuf.String())
	}
}
