// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package log

import (
	"testing"
)

var colorTests = []struct {
	escapeCodes []interface{}
	output      string
}{
	{[]interface{}{ANSI_BOLD, "bold text attribute"},
		"\x1b[1mbold text attribute\x1b[0m"},
	{[]interface{}{ANSI_UNDERLINE, "underline text attribute"},
		"\x1b[4munderline text attribute\x1b[0m"},
	{[]interface{}{ANSI_BLINK, "blink text attribute"},
		"\x1b[5mblink text attribute\x1b[0m"},
	{[]interface{}{ANSI_REVERSE, "reverse text attribute"},
		"\x1b[7mreverse text attribute\x1b[0m"},
	{[]interface{}{ANSI_CONCEALED, "concealed text attribute"},
		"\x1b[8mconcealed text attribute\x1b[0m"},
	{[]interface{}{ANSI_BLACK, "black foreground color"},
		"\x1b[30mblack foreground color\x1b[0m"},
	{[]interface{}{ANSI_RED, "red foreground color"},
		"\x1b[31mred foreground color\x1b[0m"},
	{[]interface{}{ANSI_GREEN, "green foreground color"},
		"\x1b[32mgreen foreground color\x1b[0m"},
	{[]interface{}{ANSI_YELLOW, "yellow foreground color"},
		"\x1b[33myellow foreground color\x1b[0m"},
	{[]interface{}{ANSI_BLUE, "blue foreground color"},
		"\x1b[34mblue foreground color\x1b[0m"},
	{[]interface{}{ANSI_MAGENTA, "magenta foreground color"},
		"\x1b[35mmagenta foreground color\x1b[0m"},
	{[]interface{}{ANSI_CYAN, "cyan foreground color"},
		"\x1b[36mcyan foreground color\x1b[0m"},
	{[]interface{}{ANSI_WHITE, "white foreground color"},
		"\x1b[37mwhite foreground color\x1b[0m"},
	{[]interface{}{ANSI_BG_YELLOW, ANSI_BOLD, ANSI_RED,
		"bold red text with yellow background",
		ANSI_BG_CYAN, ANSI_GREEN, "green text with cyan background"},
		"\x1b[43m\x1b[1m\x1b[31mbold red text with yellow background" +
			"\x1b[46m\x1b[32mgreen text with cyan background\x1b[0m"},
	{[]interface{}{ANSI_UNDERLINE, ANSI_GREEN, "green underline text"},
		"\x1b[4m\x1b[32mgreen underline text\x1b[0m"},
	{[]interface{}{ANSI_BOLD, ANSI_UNDERLINE, ANSI_GREEN,
		"bold green underline text"},
		"\x1b[1m\x1b[4m\x1b[32mbold green underline text\x1b[0m"},
	{[]interface{}{ANSI_BOLD, ANSI_GREEN, "colored ", ANSI_OFF,
		"to normal text"},
		"\x1b[1m\x1b[32mcolored \x1b[0mto normal text\x1b[0m"},
	{[]interface{}{ANSI_BOLD, ANSI_GREEN, "colored", ANSI_OFF},
		"\x1b[1m\x1b[32mcolored\x1b[0m"},
}

func TestColors(t *testing.T) {
	for i, v := range colorTests {
		if out := AnsiEscape(v.escapeCodes...); out != v.output {
			t.Errorf("%d. Escape(%q) = %q, want %q", i,
				v.escapeCodes, out, v.output)
		}
	}
}
