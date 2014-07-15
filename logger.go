// Copyright 2013,2014 The go-elog Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

// The go-elog package is a drop in replacement for the Go standard log package
// that provides a number of enhancements. Including colored output, logging
// levels, custom log formatting, and multiple simultaneous output streams like
// os.Stdout or a File.
package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/demizer/rgbterm"
)

// Used for string output of the logging object
var levels = [6]string{
	"LEVEL_DEBUG",
	"LEVEL_INFO",
	"LEVEL_WARNING",
	"LEVEL_ERROR",
	"LEVEL_CRITICAL",
	"LEVEL_ALL",
}

// Used to retrieve a ansi colored label of the logger
var labels = [6]string{
	// Print labels for special logging functions
	rgbterm.String("[DEBUG]", 255, 255, 255),   // White
	rgbterm.String("[INFO]", 0, 215, 95),       // Green
	rgbterm.String("[WARNING]", 255, 255, 135), // Yellow
	rgbterm.String("[ERROR]", 255, 0, 215),     // Magenta
	rgbterm.String("[CRITICAL]", 255, 0, 0),    // Red
	"", // The Print* functions do not use a label
}

type level int

// Returns the string representation of the level
func (l level) String() string { return levels[l] }

// Returns the ansi colorized label for the level
func (l level) Label() string {
	return labels[l]
}

const (
	// LEVEL_DEBUG level messages should be used for development logging
	// instead of Printf calls. When used in this manner, instead of
	// sprinkling Printf calls everywhere and then having to remove them
	// once the bug is fixed, the developer can simply change to a higher
	// logging level and the debug messages will not be sent to the output
	// stream.
	LEVEL_DEBUG level = iota

	// LEVEL_INFO level messages should be used to convey more informative
	// output than debug that could be used by a user.
	LEVEL_INFO

	// LEVEL_WARNING messages should be used to notify the user that
	// something worked, but the expected value was not the result.
	LEVEL_WARNING

	// LEVEL_ERROR messages should be used when something just did not work
	// at all.
	LEVEL_ERROR

	// LEVEL_CRITICAL messages are used when something is completely broken
	// and unrecoverable. Critical messages are usually followed by
	// os.Exit().
	LEVEL_CRITICAL

	// LEVEL_ALL level shows all messages. This is used by default for the
	// Print*() functions.
	LEVEL_ALL
)

var (
	defaultDate        = "Mon-20060102-15:04:05"
	defaultPrefix      = "::"
	defaultPrefixColor = rgbterm.String("::", 0, 255, 135) // Green
)

const (
	// These flags define which text to prefix to each log entry generated
	// by the Logger. Bits or'ed together to control what's printed.
	Ldate = 1 << iota

	// Full file name and line number: /a/b/c/d.go:23
	LlongFileName

	// Base file name and line number: d.go:23. overrides LshortFileName
	LshortFileName

	// Calling function name
	LfunctionName

	// Calling function line number
	LlineNumber

	// Use color escape sequences
	Lcolor

	// Disable ansi in file output
	LnoFileAnsi

	// Disable prefix output
	LnoPrefix

	// Show ids for functions generating output. Useful for disabling
	// specific output.
	Lid

	// initial values for the standard logger
	LstdFlags = Ldate | Lcolor | LnoFileAnsi
)

// A Logger represents an active logging object that generates lines of output
// to an io.Writer. Each logging operation makes a single call to the Writer's
// Write method. A Logger can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the Writer.
type logger struct {
	mu         sync.Mutex         // Ensures atomic writes
	buf        []byte             // For marshaling output to write
	dateFormat string             // time.RubyDate is the default format
	flags      int                // Properties of the output
	level      level              // The default level is warning
	lastId     int                // The last id level encountered
	ids        map[string]int     // ids level of the log line
	template   *template.Template // The format order of the output
	prefix     string             // Inserted into every logging output
	streams    []io.Writer        // Destination for output
}

var (
	// The default logger
	std = New(LEVEL_CRITICAL, os.Stderr)
)

// New creates a new logger object and returns it.
func New(level level, streams ...io.Writer) (obj *logger) {
	tmpl := template.Must(template.New("default").Funcs(funcMap).Parse(logFmt))
	obj = &logger{
		ids:        make(map[string]int),
		streams:    streams,
		dateFormat: defaultDate,
		flags:      LstdFlags,
		level:      level,
		template:   tmpl,
		prefix:     defaultPrefixColor,
	}
	return
}

// Returns the template of the standard logging object.
func Template() *template.Template { return std.template }

// SetTemplate allocates and parses a new output template for the logging
// object.
func SetTemplate(temp string) error {
	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		return err
	}
	std.template = tmpl
	return nil
}

// Returns the date format used by the standard logging object as a string.
func DateFormat() string { return std.dateFormat }

// Set the date format of the standard logging object. See the date package
// documentation for details on using the date format string.
func SetDateFormat(format string) { std.dateFormat = format }

// Returns the usages flags of the standard logging object.
func Flags() int { return std.flags }

// Set the usage flags for the standard logging object.
func SetFlags(flags int) { std.flags = flags }

// Get the logging level of the standard logging object.
func Level() level { return std.level }

// Set the logging level of the standard logging object.
func SetLevel(level level) { std.level = level }

// Get the logging prefix used by the standard logging object. By default it is
// "::".
func Prefix() string { return std.prefix }

// Set the logging prefix of the standard logging object.
func SetPrefix(prefix string) { std.prefix = prefix }

// Get the output streams of the standard logger
func Streams() []io.Writer { return std.streams }

// Set the output streams of the standard logger
func SetStreams(streams ...io.Writer) { std.streams = streams }

// Printf formats according to a format specifier and writes to standard
// logger output stream(s).
func Printf(format string, v ...interface{}) {
	std.Fprint(LEVEL_ALL, 2, fmt.Sprintf(format, v...), nil)
}

// Print sends output to the standard logger object output stream(s) regardless
// of logging level. The output is formatted using the output template and
// flags. Spaces are added between operands when neither is a string.
func Print(v ...interface{}) {
	std.Fprint(LEVEL_ALL, 2, fmt.Sprint(v...), nil)
}

// Println formats using the default formats for its operands and writes to the
// standard logger output stream(s). Spaces are always added between operands and
// a newline is appended.
func Println(v ...interface{}) {
	std.Fprint(LEVEL_ALL, 2, fmt.Sprintln(v...), nil)
}

// Fatalf is equivalent to Printf(), but will terminate the program with
// os.Exit(1) once output is complete.
func Fatalf(format string, v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

// Fatal is equivalent to Print(). Once output is finished, os.Exit(1) is used
// to terminate the program.
func Fatal(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	os.Exit(1)
}

// Fatalln is equivalent to Println(). Once output is finished, os.Exit(1) is
// used to terminate the program.
func Fatalln(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

// Panicf is equivalent to Printf(), but panic() is called once output is
// complete.
func Panicf(format string, v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	panic(v)
}

// Panic is equivalent to Print(), but panic() is called once output is
// complete.
func Panic(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	panic(v)
}

// Panicln is equivalent to Println(), but panic() is called once output is
// complete.
func Panicln(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	panic(v)
}

// Debugf is similar to Printf(), except the colorized LEVEL_DEBUG label is
// prefixed to the output.
func Debugf(format string, v ...interface{}) {
	std.Fprint(LEVEL_DEBUG, 2, fmt.Sprintf(format, v...), nil)
}

// Debug is similar to Print(), except the colorized LEVEL_DEBUG label is
// prefixed to the output.
func Debug(v ...interface{}) {
	std.Fprint(LEVEL_DEBUG, 2, fmt.Sprint(v...), nil)
}

// Debugln is similar to Println(), except the colorized LEVEL_DEBUG label is
// prefixed to the output.
func Debugln(v ...interface{}) {
	std.Fprint(LEVEL_DEBUG, 2, fmt.Sprintln(v...), nil)
}

// Infof is similar to Printf(), except the colorized LEVEL_INFO label is
// prefixed to the output.
func Infof(format string, v ...interface{}) {
	std.Fprint(LEVEL_INFO, 2, fmt.Sprintf(format, v...), nil)
}

// Info is similar to Print(), except the colorized LEVEL_INFO label is prefixed
// to the output.
func Info(v ...interface{}) {
	std.Fprint(LEVEL_INFO, 2, fmt.Sprint(v...), nil)
}

// Infoln is similar to Println(), except the colorized LEVEL_INFO label is
// prefixed to the output.
func Infoln(v ...interface{}) {
	std.Fprint(LEVEL_INFO, 2, fmt.Sprintln(v...), nil)
}

// Warningf is similar to Printf(), except the colorized LEVEL_WARNING label is
// prefixed to the output.
func Warningf(format string, v ...interface{}) {
	std.Fprint(LEVEL_WARNING, 2, fmt.Sprintf(format, v...), nil)
}

// Warning is similar to Print(), except the colorized LEVEL_WARNING label is
// prefixed to the output.
func Warning(v ...interface{}) {
	std.Fprint(LEVEL_WARNING, 2, fmt.Sprint(v...), nil)
}

// Warningln is similar to Println(), except the colorized LEVEL_WARNING label
// is prefixed to the output.
func Warningln(v ...interface{}) {
	std.Fprint(LEVEL_WARNING, 2, fmt.Sprintln(v...), nil)
}

// Errorf is similar to Printf(), except the colorized LEVEL_ERROR label is
// prefixed to the output.
func Errorf(format string, v ...interface{}) {
	std.Fprint(LEVEL_ERROR, 2, fmt.Sprintf(format, v...), nil)
}

// Error is similar to Print(), except the colorized LEVEL_ERROR label is
// prefixed to the output.
func Error(v ...interface{}) {
	std.Fprint(LEVEL_ERROR, 2, fmt.Sprint(v...), nil)
}

// Errorln is similar to Println(), except the colorized LEVEL_ERROR label is
// prefixed to the output.
func Errorln(v ...interface{}) {
	std.Fprint(LEVEL_ERROR, 2, fmt.Sprintln(v...), nil)
}

// Criticalf is similar to Printf(), except the colorized LEVEL_CRITICAL label is
// prefixed to the output.
func Criticalf(format string, v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
}

// Critical is similar to Prin()t, except the colorized LEVEL_CRITICAL label is
// prefixed to the output.
func Critical(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
}

// Criticalln is similar to Println(), except the colorized LEVEL_CRITICAL label
// is prefixed to the output.
func Criticalln(v ...interface{}) {
	std.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
}

// Fprint is used by all of the logging functions to send output to the output
// stream.
//
// logLevel is the level of the output.
//
// calldepth is the number of stack frames to skip when getting the file
// name of original calling function for file name output.
//
// text is the string to append to the assembled log format output. If the text
// is prefixed with newlines, they will be stripped out and placed in front of
// the completed output (test with template applied) before writing it to the
// stream.
//
// stream will be used as the output stream the text will be written to. If
// stream is nil, the stream value contained in the logger object is used.
//
// Fprint returns the number of bytes written to the stream or an error.
func (l *logger) Fprint(logLevel level, calldepth int,
	text string, stream io.Writer) (n int, err error) {

	if (logLevel != LEVEL_ALL && l.level != LEVEL_ALL) &&
		logLevel < l.level {
		return 0, nil
	}

	now := time.Now()
	var pgmC uintptr
	var file, fName string
	var line int
	var id string

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.flags&(LlongFileName|LshortFileName|LfunctionName|Lid) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		pgmC, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		if l.flags&Lid != 0 {
			fAtPC := runtime.FuncForPC(pgmC)
			hName := fAtPC.Name()
			var idNum int
			if _, ok := l.ids[hName]; ok {
				idNum = l.ids[hName]
				// fmt.Println("EXIST", hName, idNum)
			} else {
				l.ids[hName] = l.lastId
				idNum = l.lastId
				l.lastId += 1
				// fmt.Println("NEW  ", hName, idNum)
			}
			id = fmt.Sprintf("[%02.f]", float64(idNum))
		}
		if l.flags&LshortFileName != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		if l.flags&LfunctionName != 0 {
			fAtPC := runtime.FuncForPC(pgmC)
			fName = fAtPC.Name()
			for i := len(fName) - 1; i >= 0; i-- {
				if fName[i] == '.' {
					fName = fName[i+1:]
					break
				}
			}
		}
		l.mu.Lock()
	}

	// Reset the buffer
	l.buf = l.buf[:0]

	trimText := strings.TrimLeft(text, "\n")
	trimedCount := len(text) - len(trimText)
	if trimedCount > 0 {
		l.buf = append(l.buf, trimText...)
	} else {
		l.buf = append(l.buf, text...)
	}

	var date string
	var prefix string

	if l.flags&(Ldate) != 0 {
		date = now.Format(l.dateFormat)
	}

	if l.flags&(LnoPrefix) == 0 {
		prefix = l.prefix
	}

	if l.flags&(LlongFileName|LshortFileName) == 0 {
		file = ""
	}

	if l.flags&(LlineNumber) == 0 {
		line = 0
	}

	f := &format{
		Prefix:       prefix,
		LogLabel:     logLevel.Label(),
		Date:         date,
		FileName:     file,
		FunctionName: fName,
		LineNumber:   line,
		Id:           id,
		Text:         string(l.buf),
	}

	var out bytes.Buffer
	var strippedText, finalText string

	err = l.template.Execute(&out, f)

	if l.flags&Lcolor == 0 {
		strippedText = stripAnsi(out.String())
	}

	if trimedCount > 0 && l.flags&Lcolor == 0 {
		finalText = strings.Repeat("\n", trimedCount) + strippedText
	} else if trimedCount > 0 && l.flags&Lcolor != 0 {
		finalText = strings.Repeat("\n", trimedCount) + out.String()
	} else if l.flags&Lcolor == 0 {
		finalText = strippedText
	} else {
		finalText = out.String()
	}

	if stream == nil {
		n, err = l.Write([]byte(finalText))
	} else {
		n, err = stream.Write([]byte(finalText))
	}

	return
}

// Returns the template of the standard logging object.
func (l *logger) Template() *template.Template { return l.template }

// SetTemplate allocates and parses a new output template for the logging
// object.
func (l *logger) SetTemplate(temp string) error {
	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		return err
	}
	l.template = tmpl
	return nil
}

// Returns the date format used by the logging object as a string.
func (l *logger) DateFormat() string { return l.dateFormat }

// Set the date format of the logging object. See the date package
// documentation for details on using the date format string.
func (l *logger) SetDateFormat(format string) { l.dateFormat = format }

// Returns the usages flags of the logging object.
func (l *logger) Flags() int { return l.flags }

// Set the usage flags for the logging object.
func (l *logger) SetFlags(flags int) { l.flags = flags }

// Get the logging level of the logging object.
func (l *logger) Level() level { return l.level }

// Set the logging level of the logging object.
func (l *logger) SetLevel(level level) { l.level = level }

// Get the logging prefix used by the logging object. By default it is "::".
func (l *logger) Prefix() string { return l.prefix }

// Set the logging prefix of the logging object.
func (l *logger) SetPrefix(prefix string) { l.prefix = prefix }

// Get the output streams of the logger
func (l *logger) Streams() []io.Writer { return l.streams }

// Set the output streams of the logger
func (l *logger) SetStreams(streams ...io.Writer) { l.streams = streams }

// Write writes the array of bytes (p) to all of the logger.Streams. If the
// Lcolor flag is set, ansi escape codes are used to add coloring to the output.
func (l *logger) Write(p []byte) (n int, err error) {
	for _, w := range l.streams {
		if reflect.TypeOf(w).String() == "*os.File" && l.flags&LnoFileAnsi != 0 {
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

// Printf is equivalent to log.Printf().
func (l *logger) Printf(format string, v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprintf(format, v...), nil)
}

// Print is equivalent to log.Print().
func (l *logger) Print(v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprint(v...), nil)
}

// Println is equivalent to log.Println().
func (l *logger) Println(v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprintln(v...), nil)
}

// Fatalf is equivalent to log.Fatalf().
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

// Fatal is equivalent to log.Fatal().
func (l *logger) Fatal(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	os.Exit(1)
}

// Fatalln is equivalent to log.Fatalln().
func (l *logger) Fatalln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

// Panicf is equivalent to log.Panicf().
func (l *logger) Panicf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	panic(v)
}

// Panic is equivalent to log.Panic().
func (l *logger) Panic(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	panic(v)
}

// Panicln is equivalent to log.Panicln().
func (l *logger) Panicln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	panic(v)
}

// Debugf is equivalent to log.Debugf().
func (l *logger) Debugf(format string, v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprintf(format, v...), nil)
}

// Debug is equivalent to log.Debug().
func (l *logger) Debug(v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprint(v...), nil)
}

// Debugln is equivalent to log.Debugln().
func (l *logger) Debugln(v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprintln(v...), nil)
}

// Infof is equivalent to log.Infof().
func (l *logger) Infof(format string, v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprintf(format, v...), nil)
}

// Info is equivalent to log.Info().
func (l *logger) Info(v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprint(v...), nil)
}

// Infoln is equivalent to log.Infoln().
func (l *logger) Infoln(v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprintln(v...), nil)
}

// Warningf is equivalent to log.Warningf().
func (l *logger) Warningf(format string, v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprintf(format, v...), nil)
}

// Warning is equivalent to log.Warning().
func (l *logger) Warning(v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprint(v...), nil)
}

// Warningln is equivalent to log.Warningln().
func (l *logger) Warningln(v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprintln(v...), nil)
}

// Errorf is equivalent to log.Errorf().
func (l *logger) Errorf(format string, v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprintf(format, v...), nil)
}

// Error is equivalent to log.Error().
func (l *logger) Error(v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprint(v...), nil)
}

// Errorln is equivalent to log.Errorln().
func (l *logger) Errorln(v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprintln(v...), nil)
}

// Criticalf is equivalent to log.Criticalf().
func (l *logger) Criticalf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
}

// Critical is equivalent to log.Critical().
func (l *logger) Critical(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
}

// Criticalln is equivalent to log.Criticalln().
func (l *logger) Criticalln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
}
