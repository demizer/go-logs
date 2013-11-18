// Copyright 2013 The go-elog Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package log

import (
	"fmt"
)

// eCode is an ANSI escape code
type eCode int

// Ansi escape code constants. See
// http://ascii-table.com/ansi-escape-sequences.php

const (
	// General text attributes
	ANSI_OFF  eCode = iota
	ANSI_BOLD       // 1
	_
	_
	ANSI_UNDERLINE // 4
	ANSI_BLINK     // 5
	_
	ANSI_REVERSE   // 7
	ANSI_CONCEALED // 8
)

const (
	// Foreground text attributes
	ANSI_BLACK   eCode = iota + 30
	ANSI_RED           // 31
	ANSI_GREEN         // 32
	ANSI_YELLOW        // 33
	ANSI_BLUE          // 34
	ANSI_MAGENTA       // 35
	ANSI_CYAN          // 36
	ANSI_WHITE         // 37
)

const (
	// Background text attributes
	ANSI_BG_GREY    eCode = iota + 40
	ANSI_BG_RED           // 41
	ANSI_BG_GREEN         // 42
	ANSI_BG_YELLOW        // 43
	ANSI_BG_BLUE          // 44
	ANSI_BG_MAGENTA       // 45
	ANSI_BG_CYAN          // 46
	ANSI_BG_WHITE         // 47
)

// AnsiEscape accepts ANSI escape codes and strings to form escape sequences.
// For example, to create a string with a colorized prefix,
//
//      AnsiEscape(ANSI_BOLD, ANSI_GREEN, "[DEBUG] ", ANSI_OFF, "Text string")
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
	if c[len(c)-1] != ANSI_OFF {
		out += "\x1b[0m"
	}
	return
}
