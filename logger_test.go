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
	eLen := 55
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
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
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
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

var (
	boldPrefix  = AnsiEscape(ANSI_BOLD, "TEST>", ANSI_OFF)
	colorPrefix = AnsiEscape(ANSI_BOLD, ANSI_RED, "TEST>", ANSI_OFF)
	date        = "Mon 20060102 15:04:05"
)

var outputTests = []struct {
	template   string
	prefix     string
	level      level
	dateFormat string
	flags      int
	text       string
	want       string
	wantErr    bool
}{

	// The %s format specifier is the placeholder for the date.
	{logFmt, boldPrefix, LEVEL_ALL, date, LstdFlags, "test number 1",
		"%s \x1b[1mTEST>\x1b[0m test number 1", false},

	{logFmt, colorPrefix, LEVEL_ALL, date, LstdFlags, "test number 2",
		"%s \x1b[1m\x1b[31mTEST>\x1b[0m test number 2", false},

	// Test output with coloring turned off
	{logFmt, AnsiEscape(ANSI_BOLD, "::", ANSI_OFF), LEVEL_ALL, date, Ldate,
		"test number 3", "%s :: test number 3", false},

	{logFmt, defaultPrefixColor, LEVEL_DEBUG, time.RubyDate, LstdFlags,
		"test number 4",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[37m[DEBUG]\x1b[0m test number 4",
		false},

	{logFmt, defaultPrefixColor, LEVEL_INFO, time.RubyDate, LstdFlags,
		"test number 5",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[32m[INFO]\x1b[0m test number 5",
		false},

	{logFmt, defaultPrefixColor, LEVEL_WARNING, time.RubyDate, LstdFlags,
		"test number 6",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[33m[WARNING]\x1b[0m test number 6",
		false},

	{logFmt, defaultPrefixColor, LEVEL_ERROR, time.RubyDate, LstdFlags,
		"test number 7",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[35m[ERROR]\x1b[0m test number 7",
		false},

	{logFmt, defaultPrefixColor, LEVEL_CRITICAL, time.RubyDate, LstdFlags,
		"test number 8",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[31m[CRITICAL]\x1b[0m test number 8",
		false},

	// Test date format
	{logFmt, defaultPrefixColor, LEVEL_ALL, "Mon 20060102 15:04:05",
		Ldate, "test number 9",
		"%s :: test number 9", false},
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
		want := fmt.Sprintf(k.want, d)
		if buf.String() != want || err != nil && !k.wantErr {
			t.Errorf("Print test %d failed, \ngot:  %q\nwant: "+
				"%q", i+1, buf.String(), want)
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
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix)
	Debug("\n\nThis line should be padded with newlines.\n\n")
	expect := "\n\n[DEBUG] This line should be padded with newlines.\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n%q\nGot:\n%q\n", expect, buf.String())
	}
}

func TestFlagsLdate(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix)
	Debugln("This output should not have a date.")
	expect := "[DEBUG] This output should not have a date.\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsLfunctionName(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | LfunctionName)
	Debugln("This output should have a function name.")
	expect := "[DEBUG] TestFlagsLfunctionName: This output should have a function name.\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsLfunctionNameWithFileName(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | LfunctionName | LshortFileName)
	Debug("This output should have a file name and a function name.")
	expect := "[DEBUG] logger_test.go: TestFlagsLfunctionNameWithFileName" +
		": This output should have a file name and a function name."
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsLansiWithNewlinePaddingDebug(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | Lansi)
	Debug("\n\nThis output should be padded with newlines and colored.\n\n")
	expect := "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m This output should be " +
		"padded with newlines and colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsLansiWithNewlinePaddingDebugf(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | Lansi)
	Debugf("\n\nThis output should be padded with newlines and %s.\n\n",
		"colored")
	expect := "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m This output should be " +
		"padded with newlines and colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
	buf.Reset()
	Debugf("\n\n##### HELLO %s #####\n\n", "NEWMAN")
	expect = "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m ##### HELLO NEWMAN #####\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsLansiWithNewlinePaddingDebugln(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix | Lansi)
	Debugln("\n\nThis output should be padded with newlines and colored.\n\n")
	expect := "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m This output should be " +
		"padded with newlines and colored.\n\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
	buf.Reset()
	Debugln("\n\n", "### HELLO", "NEWMAN", "###", "\n\n")
	expect = "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m  ### HELLO NEWMAN ### \n\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
	buf.Reset()
	Debugln("\n\n### HELLO", "NEWMAN", "###\n\n")
	expect = "\n\n\x1b[1m\x1b[37m[DEBUG]\x1b[0m ### HELLO NEWMAN ###\n\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestFlagsNoLansiWithNewlinePadding(t *testing.T) {
	var buf bytes.Buffer
	SetStreams(&buf)
	SetLevel(LEVEL_DEBUG)
	SetFlags(LnoPrefix)
	Debug("\n\nThis output should be padded with newlines and not colored.\n\n")
	expect := "\n\n[DEBUG] This output should be padded with newlines and not colored.\n\n"
	if buf.String() != expect {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", expect, buf.String())
	}
}

func TestHeirarchicalPrintln(t *testing.T) {
	var buf bytes.Buffer
	var tBuf bytes.Buffer

	SetStreams(&buf)
	SetLevel(LEVEL_WARNING)
	SetFlags(LstdFlags | Lheirarchical)

	now := time.Now()

	Println("\n\nLevel 0 Output 1")
	lvl2 := func() {
		Println("Level 2 Output 1")
		Println("Level 2 Output 2")
	}
	lvl1 := func() {
		Println("Level 1 Output 1")
		Println("Level 1 Output 2")
		lvl2()
	}
	lvl1()

	date = now.Format(std.DateFormat)
	f := struct {
		Date string
	}{Date: date}

	temp := "\n\n{{.Date}} :: [00] Level 0 Output 1\n" +
		"{{.Date}} ::      [01] Level 1 Output 1\n" +
		"{{.Date}} ::      [01] Level 1 Output 2\n" +
		"{{.Date}} ::           [02] Level 2 Output 1\n" +
		"{{.Date}} ::           [02] Level 2 Output 2\n"

	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		t.Fatal(err)
	}
	tmpl.Execute(&tBuf, f)

	if buf.String() != tBuf.String() {
		t.Errorf("\nExpect:\n\t%q\nGot:\n\t%q\n", tBuf.String(), buf.String())
	}
}
