// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package logger

import (
	"fmt"
)

// eCode is an ANSI escape code
type eCode int

// Ansi escape code constants. See
// http://ascii-table.com/ansi-escape-sequences.php

const (
	// General text attributes
	OFF  eCode = iota
	BOLD       // 1
	_
	_
	UNDERLINE // 4
	BLINK     // 5
	_
	REVERSE   // 7
	CONCEALED // 8
)

const (
	// Foreground text attributes
	BLACK   eCode = iota + 30
	RED           // 31
	GREEN         // 32
	YELLOW        // 33
	BLUE          // 34
	MAGENTA       // 35
	CYAN          // 36
	WHITE         // 37
)

const (
	// Background text attributes
	BG_GREY    eCode = iota + 40
	BG_RED           // 41
	BG_GREEN         // 42
	BG_YELLOW        // 43
	BG_BLUE          // 44
	BG_MAGENTA       // 45
	BG_CYAN          // 46
	BG_WHITE         // 47
)

// AnsiEscape accepts ANSI escape codes and strings to form escape sequences.
// For example, to create a string with a colorized prefix,
//
//      AnsiEscape(BOLD, GREEN, "[DEBUG] ", OFF, "Here is the debug output")
//
// and a nicely escaped string for terminal output will be returned.
func AnsiEscape(c ...interface{}) (out string) {
	for _, val := range c {
		switch t := val.(type) {
		case eCode:
			out += fmt.Sprintf("\x1b[%dm", val)
		case string:
			out += fmt.Sprintf("%s", val)
		default:
			fmt.Printf("unexpected type: %T\n", t)
		}
	}
	if c[len(c)-1] != OFF {
		out += "\x1b[0m"
	}
	return
}
