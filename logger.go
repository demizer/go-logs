// Copyright 2013 The go-elog Authors. All rights reserved.
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
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
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
	AnsiEscape(ANSI_BOLD, ANSI_WHITE, "[DEBUG]", ANSI_OFF),
	AnsiEscape(ANSI_BOLD, ANSI_GREEN, "[INFO]", ANSI_OFF),
	AnsiEscape(ANSI_BOLD, ANSI_YELLOW, "[WARNING]", ANSI_OFF),
	AnsiEscape(ANSI_BOLD, ANSI_MAGENTA, "[ERROR]", ANSI_OFF),
	AnsiEscape(ANSI_BOLD, ANSI_RED, "[CRITICAL]", ANSI_OFF),
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
	defaultPrefixColor = AnsiEscape(ANSI_BOLD, ANSI_GREEN, "::", ANSI_OFF)
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

	// Use ansi escape sequences
	Lansi

	// Disable ansi in file output
	LnoFileAnsi

	// Disable prefix output
	LnoPrefix

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

var (
	// The default logger
	std = New(LEVEL_CRITICAL, os.Stderr)
)

// New creates a new logger object and returns it.
func New(level level, streams ...io.Writer) (obj *Logger) {
	tmpl := template.Must(template.New("default").Funcs(funcMap).Parse(logFmt))
	obj = &Logger{Streams: streams, DateFormat: defaultDate,
		Flags: LstdFlags, Level: level, Template: tmpl, Prefix: defaultPrefixColor}
	return
}

// Returns the template of the standard logging object.
func Template() *template.Template { return std.Template }

// SetTemplate allocates and parses a new output template for the logging
// object.
func SetTemplate(temp string) error {
	tmpl, err := template.New("default").Funcs(funcMap).Parse(temp)
	if err != nil {
		return err
	}
	std.Template = tmpl
	return nil
}

// Returns the date format used by the standard logging object as a string.
func DateFormat() string { return std.DateFormat }

// Set the date format of the standard logging object. See the date package
// documentation for details on using the date format string.
func SetDateFormat(format string) { std.DateFormat = format }

// Returns the usages flags of the standard logging object.
func Flags() int { return std.Flags }

// Set the usage flags for the standard logging object.
func SetFlags(flags int) { std.Flags = flags }

// Get the logging level of the standard logging object.
func Level() level { return std.Level }

// Set the logging level of the standard logging object.
func SetLevel(level level) { std.Level = level }

// Get the logging prefix used by the standard logging object. By default it is
// "::".
func Prefix() string { return std.Prefix }

// Set the logging prefix of the standard logging object.
func SetPrefix(prefix string) { std.Prefix = prefix }

// Get the output streams of the standard logger
func Streams() []io.Writer { return std.Streams }

// Set the output streams of the standard logger
func SetStreams(streams ...io.Writer) { std.Streams = streams }

// Printf formats according to a format specifier and writes to standard
// logger output stream(s).
func Printf(format string, v ...interface{}) {
	std.Fprint(LEVEL_ALL, 2, fmt.Sprintf(format, v...), nil)
}

// Print sends output to the standard logger object output stream(s) regardless
// of logging level. The output is formatted using the output template and
// flags. Spaces are added between operands when neither is a string.
func Print(v ...interface{}) {
	std.Fprint(LEVEL_ALL, 1, fmt.Sprint(v...), nil)
}

// Println formats using the default formats for its operands and writes to the
// standard logger output stream(s). Spaces are always added between operands and
// a newline is appended.
func Println(v ...interface{}) {
	std.Fprint(LEVEL_ALL, 2, fmt.Sprintln(v...), nil)
}

//
// TODO: Need to test this!!
//
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

//
// TODO: Need to test this!!
//
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
func (l *Logger) Fprint(logLevel level, calldepth int,
	text string, stream io.Writer) (n int, err error) {

	if (logLevel != LEVEL_ALL && l.Level != LEVEL_ALL) &&
		logLevel < l.Level {
		return 0, nil
	}

	now := time.Now()
	var pgmC uintptr
	var file,fName string
	var line int

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Flags&(LlongFileName|LshortFileName|LfunctionName) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		pgmC, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		if l.Flags&LshortFileName != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		if l.Flags&LfunctionName != 0 {
			fAtPC := runtime.FuncForPC(pgmC)
			fName = fAtPC.Name()
			// fmt.Println(fName)
			// fmt.Println("fname length:", len(fName))
			for i := len(fName) - 1; i >= 0; i-- {
				// fmt.Printf("i = %d, %s == %s = %+v\n", i,
					// string(fName[i]), string('.'), fName[i] == '.')
				// if file[i] == ':' {
					// endFname = i
				// }
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

	if l.Flags&(Ldate) != 0 {
		date = now.Format(l.DateFormat)
	}

	if l.Flags&(LnoPrefix) == 0 {
		prefix = l.Prefix
	}

	if l.Flags&(LlongFileName|LshortFileName) == 0 {
		file = ""
	}

	if l.Flags&(LlineNumber) == 0 {
		line = 0
	}

	f := &format{
		Prefix: prefix,
		LogLabel: logLevel.Label(),
		Date: date,
		FileName: file,
		FunctionName: fName,
		LineNumber: line,
		Text: string(l.buf),
	}

	var out bytes.Buffer
	err = l.Template.Execute(&out, f)

	if trimedCount > 0 {
		text = strings.Repeat("\n", trimedCount) + out.String()
	} else {
		text = out.String()
	}

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

// Printf is equivalent to log.Printf().
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprintf(format, v...), nil)
}

// Print is equivalent to log.Print().
func (l *Logger) Print(v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprint(v...), nil)
}

// Println is equivalent to log.Println().
func (l *Logger) Println(v ...interface{}) {
	l.Fprint(LEVEL_ALL, 2, fmt.Sprintln(v...), nil)
}

// Fatalf is equivalent to log.Fatalf().
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

// Fatal is equivalent to log.Fatal().
func (l *Logger) Fatal(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	os.Exit(1)
}

// Fatalln is equivalent to log.Fatalln().
func (l *Logger) Fatalln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

// Panicf is equivalent to log.Panicf().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
	panic(v)
}

// Panic is equivalent to log.Panic().
func (l *Logger) Panic(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
	panic(v)
}

// Panicln is equivalent to log.Panicln().
func (l *Logger) Panicln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
	panic(v)
}

// Debugf is equivalent to log.Debugf().
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprintf(format, v...), nil)
}

// Debug is equivalent to log.Debug().
func (l *Logger) Debug(v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprint(v...), nil)
}

// Debugln is equivalent to log.Debugln().
func (l *Logger) Debugln(v ...interface{}) {
	l.Fprint(LEVEL_DEBUG, 2, fmt.Sprintln(v...), nil)
}

// Infof is equivalent to log.Infof().
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprintf(format, v...), nil)
}

// Info is equivalent to log.Info().
func (l *Logger) Info(v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprint(v...), nil)
}

// Infoln is equivalent to log.Infoln().
func (l *Logger) Infoln(v ...interface{}) {
	l.Fprint(LEVEL_INFO, 2, fmt.Sprintln(v...), nil)
}

// Warningf is equivalent to log.Warningf().
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprintf(format, v...), nil)
}

// Warning is equivalent to log.Warning().
func (l *Logger) Warning(v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprint(v...), nil)
}

// Warningln is equivalent to log.Warningln().
func (l *Logger) Warningln(v ...interface{}) {
	l.Fprint(LEVEL_WARNING, 2, fmt.Sprintln(v...), nil)
}

// Errorf is equivalent to log.Errorf().
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprintf(format, v...), nil)
}

// Error is equivalent to log.Error().
func (l *Logger) Error(v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprint(v...), nil)
}

// Errorln is equivalent to log.Errorln().
func (l *Logger) Errorln(v ...interface{}) {
	l.Fprint(LEVEL_ERROR, 2, fmt.Sprintln(v...), nil)
}

// Criticalf is equivalent to log.Criticalf().
func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintf(format, v...), nil)
}

// Critical is equivalent to log.Critical().
func (l *Logger) Critical(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprint(v...), nil)
}

// Criticalln is equivalent to log.Criticalln().
func (l *Logger) Criticalln(v ...interface{}) {
	l.Fprint(LEVEL_CRITICAL, 2, fmt.Sprintln(v...), nil)
}
