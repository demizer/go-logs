package log

import (
	"bytes"
	"strings"
	"strconv"
	"runtime"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	if out := Streams()[0]; out != os.Stderr {
		t.Errorf("log.Stream is not stderr by default")
	}
	var buf bytes.Buffer
	Streams()[0] = &buf
	if out := Streams()[0]; out != &buf {
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
	eLen := 68
	log := New(DEBUG, file, &buf)
	log.Debugln("Testing debug output!")
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
	log := New(DEBUG, b)
	log.Flags = LstdFlags | LlongFile
	log.Debugln("testing long file flag")
	_, file, lNum, _ := runtime.Caller(0)
	dOut := b.String()
	if strings.Index(dOut, file) < 0 {
		t.Errorf("Debugln() = %q; does not contain %s", dOut, file)
	}
	lSrch := ".go:" + strconv.Itoa(lNum - 1)
	if strings.Index(dOut, lSrch) < 0 {
		t.Errorf("Debugln() = %q; does not contain %q", dOut, lSrch)
	}
}

func TestShortFileFlag(t *testing.T) {
	b := new(bytes.Buffer)
	log := New(DEBUG, b)
	log.Flags = LstdFlags | LshortFile
	log.Debugln("testing short file flag")
	_, file, lNum, _ := runtime.Caller(0)
	sName := filepath.Base(file)
	dOut := b.String()
	if strings.Index(dOut, sName) < 0 || strings.Index(dOut, file) > 0 {
		t.Errorf("Debugln() = %q; does not contain %s", dOut, file)
	}
	lSrch := ".go:" + strconv.Itoa(lNum - 1)
	if strings.Index(dOut, lSrch) < 0 {
		t.Errorf("Debugln() = %q; does not contain %q", dOut, lSrch)
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
		log := New(DEBUG, &buf)
		log.Prefix = k.prefix
		log.DateFormat = k.dateFormat
		log.Flags = k.flags
		d := time.Now().Format(log.DateFormat)
		n, err := log.Fprint(k.logPrefix, 1, k.text, &buf)
		if n != buf.Len() {
			t.Error("Error: ", io.ErrShortWrite)
		}
		want := fmt.Sprintf(k.want, d)
		if buf.String() != want || err != nil && !k.wantErr {
			t.Errorf("Print test %d failed, \ngot:  %q\nwant: "+
				"%q", i, buf.String(), want)
			continue
		}
		fmt.Printf("Test %d OK: %s\n", i, buf.String())
	}
}
