// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

// Package logger is a better logging system for Go than the generic log
// package in the Go Standard Library. The logger packages provides colored
// output, logging levels, and simultaneous logging output to stdout, stderr,
// and file streams.
package logger

import (
	"log"
	"os"
)

// Used for string output of the logging object
var levels = [5]string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

type LoggingLevel int

// Returns the string representation of the level
func (l LoggingLevel) String() string { return levels[l] }

// The DEBUG level is the lowest possible output level. This is meant for
// development use. The default output level is WARNING.
const (
	DEBUG    LoggingLevel = iota // Used for development
	INFO                         // Used to indicate extra information
	WARNING                      // Indicates something is wrong, but not broken
	ERROR                        // Something is broken
	CRITICAL                     // A message so bad it ends the application process
)

var (
	Colors = true // Enable/Disable colored output

	// The default level only displays WARNING messages and higher.
	Level = WARNING
)

// StdLogger is the default logger. By default it directs output to stderr.
var StdLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	StdLogger = l
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	if Level <= DEBUG {
		if Colors {
			// StdLogger.Printf("%s %v\n", e_DEBUG, v)
		} else {
			StdLogger.Printf("[D] %v\n", v)
		}
	}
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	if Level <= INFO {
		StdLogger.Printf("[I] %v\n", v)
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	if Level <= WARNING {
		StdLogger.Printf("[W] %v\n", v)
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	if Level <= ERROR {
		StdLogger.Printf("[E] %v\n", v)
	}
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	if Level <= CRITICAL {
		StdLogger.Printf("[C] %v\n", v)
	}
}
