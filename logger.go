// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

// Package logger is a better logging system for Go than the generic log
// package in the Go Standard Library. The logger packages provides colored
// output, logging levels, custom log formatting, and multiple simultaneous
// output streams like stdout or a file.
package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"text/template"
	"time"
)

// Used for string output of the logging object
var levels = [5]string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

type level int

// Returns the string representation of the level
func (l level) String() string { return levels[l] }

const (
	// DEBUG level messages should be used for development logging instead
	// of Printf calls. When used in this manner, instead of sprinkling
	// Printf calls everywhere and then having to remove them once the bug
	// is fixed, the developer can simply change to a higher logging level
	// and the debug messages will not be sent to the output stream.
	DEBUG level = iota
	// Info level messages should be used to convey more informative output
	// than debug that could be used by a user.
	INFO
	// Warning messages should be used to notify the user that something
	// worked, but the expected value was not the result.
	WARNING
	// Error messages should be used when something just did not work at
	// all.
	ERROR
	// Critical messages are used when something is completely broken and
	// unrecoverable. Critical messages are usually followed by os.Exit().
	CRITICAL
)

// logLabel is a string that is added to the output depending on the log
// function used.
type logLabel string

var (
	defaultDate        = "Mon-20060102-15:04:05"
	defaultPrefix      = "::"
	defaultPrefixColor = AnsiEscape(BOLD, GREEN, "::", OFF)

	// Print labels for special logging functions
	PrintLabel    = logLabel("")
	DebugLabel    = logLabel(AnsiEscape(BOLD, WHITE, "[DEBUG]", OFF))
	InfoLabel     = logLabel(AnsiEscape(BOLD, GREEN, "[INFO]", OFF))
	WarningLabel  = logLabel(AnsiEscape(BOLD, YELLOW, "[WARNING]", OFF))
	ErrorLabel    = logLabel(AnsiEscape(BOLD, MAGENTA, "[ERROR]", OFF))
	CriticalLabel = logLabel(AnsiEscape(BOLD, RED, "[CRITICAL]", OFF))
)

const (
	// These flags define which text to prefix to each log entry generated
	// by the Logger. Bits or'ed together to control what's printed.
	Ldate = 1 << iota
	// full file name and line number: /a/b/c/d.go:23
	LlongFile
	// base file name and line number: d.go:23. overrides Llongfile
	LshortFile
	// Use ansi escape sequences
	Lansi
	// Disable ansi in file output
	LnoFileAnsi
	// initial values for the standard logger
	LstdFlags = Ldate | Lansi | LnoFileAnsi
)

// A Logger represents an active logging object that generates lines of output
// to an io.Writer. Each logging operation makes a single call to the Writer's
// Write method. A Logger can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the Writer.
type Logger struct {
	mu         sync.Mutex         // Ensures atomic writes
	buf        []byte             // For marshaling output to write
	DateFormat string             // time.RubyDate is the default format
	Flags      int                // Properties of the output
	Level      level              // The default level is warning
	Template   *template.Template // The format order of the output
	Prefix     string             // Inserted into every logging output
	Streams    []io.Writer        // Destination for output
}

// New creates a new logger object and returns it.
func New(level level, streams ...io.Writer) (obj *Logger) {
	tmpl := template.Must(template.New("default").Funcs(funcMap).Parse(logFmt))
	obj = &Logger{Streams: streams, DateFormat: defaultDate,
		Flags: LstdFlags, Level: level, Template: tmpl, Prefix: defaultPrefixColor}
	return
}

// SetTemplate allocates and parses a new output template for the logging
// object.
func (l *Logger) SetTemplate(temp string) error {
	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		return err
	}
	l.Template = tmpl
	return nil
}

// Write writes the array of bytes (p) to all of the logger.Streams. If the
// Lansi flag is set, ansi escape codes are used to add coloring to the output.
func (l *Logger) Write(p []byte) (n int, err error) {
	for _, w := range l.Streams {
		if w != os.Stdout && w != os.Stderr && w != os.Stdin &&
			l.Flags&LnoFileAnsi != 0 {
			p = stripAnsiByte(p)
			n, err = w.Write(p)
		} else {
			n, err = w.Write(p)
		}
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// Fprint is used by all of the logging functions to send output to the output
// stream.
//
// logLabel is the prefix that should be included with the output.
//
// calldepth is the number of stack frames to skip when getting the file
// name of original calling function for file name output.
//
// text is the string to append to the assembled log format output.
//
// stream will be used as the output stream the text will be written to. If
// stream is nil, the stream value contained in the logger object is used.
func (l *Logger) Fprint(logLabel logLabel, calldepth int,
	text string, stream io.Writer) (n int, err error) {
	now := time.Now()
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Flags&(LshortFile|LlongFile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		if l.Flags&LshortFile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.buf = append(l.buf, text...)
	date := now.Format(l.DateFormat)
	f := &format{l.Prefix, logLabel, date, file, line, string(l.buf)}
	var out bytes.Buffer
	err = l.Template.Execute(&out, f)
	text = out.String()
	if l.Flags&Lansi == 0 {
		text = stripAnsi(out.String())
	}
	if stream == nil {
		n, err = l.Write([]byte(text))
	} else {
		n, err = stream.Write([]byte(text))
	}
	return
}

// Print sends output to the logger object output stream regardless of logging
// level including the logger format properties and flags. Spaces are added
// between operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Print(v ...interface{}) {
	l.Fprint(PrintLabel, 2, fmt.Sprint(v...), nil)
}

// Println formats using the default formats for its operands and writes to the
// output streams. Spaces are always added between operands and a newline is
// appended.
func (l *Logger) Println(v ...interface{}) {
	l.Fprint(PrintLabel, 2, fmt.Sprintln(v...), nil)
}

// Printf formats according to a format specifier and writes to standard
// output.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Fprint(PrintLabel, 2, fmt.Sprintf(format, v...), nil)
}

// Debug has the same properties as Print except the DEBUG logLabel is
// included with the output.
func (l *Logger) Debug(v ...interface{}) {
	l.Fprint(DebugLabel, 2, fmt.Sprint(v...), nil)
}

// Debugln has the same properties as Println, except the DEBUG logLabel is
// included with the output.
func (l *Logger) Debugln(v ...interface{}) {
	l.Fprint(DebugLabel, 2, fmt.Sprintln(v...), nil)
}

// Debugf has the same properties as Printf, except the DEBUG logLabel is
// included with the output.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Fprint(DebugLabel, 2, fmt.Sprintf(format, v...), nil)
}

// Info has the same properties as Print except the INFO logLabel is included
// with the output.
func (l *Logger) Info(v ...interface{}) {
	l.Fprint(InfoLabel, 2, fmt.Sprint(v...), nil)
}

// Infoln has the same properties as Println, except the INFO logLabel is
// included with the output.
func (l *Logger) Infoln(v ...interface{}) {
	l.Fprint(InfoLabel, 2, fmt.Sprintln(v...), nil)
}

// Infof has the same properties as Println, except the INFO logLabel is
// included with the output.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Fprint(InfoLabel, 2, fmt.Sprintf(format, v...), nil)
}

// Warning has the same properties as Print except the WARNING logLabel is
// included with the output.
func (l *Logger) Warning(v ...interface{}) {
	l.Fprint(WarningLabel, 2, fmt.Sprint(v...), nil)
}

// Warningln has the same properties as Println, except the WARNING logLabel
// is included with the output.
func (l *Logger) Warningln(v ...interface{}) {
	l.Fprint(WarningLabel, 2, fmt.Sprintln(v...), nil)
}

// Warningf has the same properties as Println, except the WARNING logLabel is
// included with the output.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Fprint(WarningLabel, 2, fmt.Sprintf(format, v...), nil)
}

// Error has the same properties as Print except the ERROR logLabel is
// included with the output.
func (l *Logger) Error(v ...interface{}) {
	l.Fprint(ErrorLabel, 2, fmt.Sprint(v...), nil)
}

// Errorln has the same properties as Println, except the ERROR logLabel is
// included with the output.
func (l *Logger) Errorln(v ...interface{}) {
	l.Fprint(ErrorLabel, 2, fmt.Sprintln(v...), nil)
}

// Errorf has the same properties as Println, except the ERROR logLabel is
// included with the output.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Fprint(ErrorLabel, 2, fmt.Sprintf(format, v...), nil)
}

// Critical has the same properties as Print except the CRITICAL logLabel is
// included with the output.
func (l *Logger) Critical(v ...interface{}) {
	l.Fprint(CriticalLabel, 2, fmt.Sprint(v...), nil)
}

// Criticalln has the same properties as Println, except the CRITICAL logLabel
// is included with the output.
func (l *Logger) Criticalln(v ...interface{}) {
	l.Fprint(CriticalLabel, 2, fmt.Sprintln(v...), nil)
}

// Criticalf has the same properties as Println, except the CRITICAL logLabel
// is included with the output.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.Fprint(CriticalLabel, 2, fmt.Sprintf(format, v...), nil)
}
