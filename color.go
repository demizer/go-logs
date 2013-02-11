// Copyright 2013 The go-logger Authors. All rights reserved.
// Thi code is MIT licensed. See the LICENSE file for more info.
package logger

import (
	"fmt"
)

type eCode int

// Ansi escape code constants. See
// http://ascii-table.com/ansi-escape-sequences.php
const (
	OFF eCode = iota
	BOLD
	_
	_
	UNDERLINE
	BLINK
	_
	REVERSE
	CONCEALED
)

const (
	BLACK eCode = iota + 30
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

const (
	BG_GREY eCode = iota + 40
	BG_RED
	BG_GREEN
	BG_YELLOW
	BG_BLUE
	BG_MAGENTA
	BG_CYAN
	BG_WHITE
)

/* AnsiEscape accepts ANSI escape codes and strings to form escape sequences.

   For example, to create a string with a colorized prefix,

	AnsiEscape(BOLD, GREEN, "[DEBUG] ", OFF, "Here is the debug output")

   and a nicely escaped string for terminal output will be returned. */
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
	out += "\x1b[0m"
	return
}
