// Copyright 2013 The go-logger Authors. All rights reserved.
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
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	var buf bytes.Buffer
	logr := New(CRITICAL, os.Stdout, &buf)
	logr.Streams[1] = &buf
	if out := logr.Streams[1]; out != &buf {
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
	logr := New(DEBUG, file, &buf)
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
	b := new(bytes.Buffer)
	logr := New(DEBUG, b)
	logr.Flags = LstdFlags | LlongFile
	logr.Debugln("testing long file flag")
	_, file, lNum, _ := runtime.Caller(0)
	dOut := b.String()
	if strings.Index(dOut, file) < 0 {
		t.Errorf("Debugln() = %q; does not contain %s", dOut, file)
	}
	lSrch := ".go:" + strconv.Itoa(lNum-1)
	if strings.Index(dOut, lSrch) < 0 {
		t.Errorf("Debugln() = %q; does not contain %q", dOut, lSrch)
	}
}

func TestShortFileFlag(t *testing.T) {
	b := new(bytes.Buffer)
	logr := New(DEBUG, b)
	logr.Flags = LstdFlags | LshortFile
	logr.Debugln("testing short file flag")
	_, file, lNum, _ := runtime.Caller(0)
	sName := filepath.Base(file)
	dOut := b.String()
	if strings.Index(dOut, sName) < 0 || strings.Index(dOut, file) > 0 {
		t.Errorf("Debugln() = %q; does not contain %s", dOut, file)
	}
	lSrch := ".go:" + strconv.Itoa(lNum-1)
	if strings.Index(dOut, lSrch) < 0 {
		t.Errorf("Debugln() = %q; does not contain %q", dOut, lSrch)
	}
}

var (
	boldPrefix  = AnsiEscape(BOLD, "TEST>", OFF)
	colorPrefix = AnsiEscape(BOLD, RED, "TEST>", OFF)
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
	{logFmt, boldPrefix, ALL, date, LstdFlags, "test number 1",
		"%s \x1b[1mTEST>\x1b[0m test number 1", false},

	{logFmt, colorPrefix, ALL, date, LstdFlags, "test number 2",
		"%s \x1b[1m\x1b[31mTEST>\x1b[0m test number 2", false},

	// Test output with coloring turned off
	{logFmt, AnsiEscape(BOLD, "::", OFF), ALL, date, Ldate,
		"test number 3", "%s :: test number 3", false},

	{logFmt, defaultPrefixColor, DEBUG, time.RubyDate, LstdFlags,
		"test number 4",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[37m[DEBUG]\x1b[0m test number 4",
		false},

	{logFmt, defaultPrefixColor, INFO, time.RubyDate, LstdFlags,
		"test number 5",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[32m[INFO]\x1b[0m test number 5",
		false},

	{logFmt, defaultPrefixColor, WARNING, time.RubyDate, LstdFlags,
		"test number 6",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[33m[WARNING]\x1b[0m test number 6",
		false},

	{logFmt, defaultPrefixColor, ERROR, time.RubyDate, LstdFlags,
		"test number 7",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[35m[ERROR]\x1b[0m test number 7",
		false},

	{logFmt, defaultPrefixColor, CRITICAL, time.RubyDate, LstdFlags,
		"test number 8",
		"%s \x1b[1m\x1b[32m::\x1b[0m \x1b[1m\x1b[31m[CRITICAL]\x1b[0m test number 8",
		false},

	// Test date format
	{logFmt, defaultPrefixColor, ALL, "Mon 20060102 15:04:05",
		Ldate, "test number 9",
		"%s :: test number 9", false},
}

func TestOutput(t *testing.T) {
	for i, k := range outputTests {
		var buf bytes.Buffer
		logr := New(DEBUG, &buf)
		logr.Prefix = k.prefix
		logr.DateFormat = k.dateFormat
		logr.Flags = k.flags
		logr.Level = k.level
		d := time.Now().Format(logr.DateFormat)
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
		fmt.Printf("Test %d OK: %s\n", i, buf.String())
	}
}

func TestLevel(t *testing.T) {
	var buf bytes.Buffer
	logr := New(CRITICAL, &buf)
	logr.Debug("This level should produce no output")
	if buf.Len() != 0 {
		t.Errorf("Debug() produced output at CRITICAL logging level")
	}
	logr.Level = DEBUG
	logr.Debug("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the DEBUG logging level")
	}
	buf.Reset()
	logr.Level = CRITICAL
	logr.Println("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the ALL logging level")
	}
	buf.Reset()
	logr.Level = ALL
	logr.Debug("This level should produce output")
	if buf.Len() == 0 {
		t.Errorf("Debug() did not produce output at the ALL logging level")
	}
}
