// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package log

import (
	"fmt"
	"testing"

	"github.com/demizer/rgbterm"
)

var colorTests = []struct {
	escapeCodes string
	output      string
}{
	{rgbterm.String("red foreground color", 255, 0, 0),
		"\x1b[38;5;196mred foreground color\x1b[0;00m"},
	{rgbterm.String("green foreground color", 0, 255, 0),
		"\x1b[38;5;46mgreen foreground color\x1b[0;00m"},
	{rgbterm.String("blue foreground color", 0, 0, 255),
		"\x1b[38;5;21mblue foreground color\x1b[0;00m"},
	// {[]interface{}{ANSI_BOLD, "bold text attribute"},
	// "\x1b[1mbold text attribute\x1b[0m"},
	// {[]interface{}{ANSI_UNDERLINE, "underline text attribute"},
	// "\x1b[4munderline text attribute\x1b[0m"},
	// {[]interface{}{ANSI_BLINK, "blink text attribute"},
	// "\x1b[5mblink text attribute\x1b[0m"},
	// {[]interface{}{ANSI_REVERSE, "reverse text attribute"},
	// "\x1b[7mreverse text attribute\x1b[0m"},
	// {[]interface{}{ANSI_CONCEALED, "concealed text attribute"},
	// "\x1b[8mconcealed text attribute\x1b[0m"},
	// {[]interface{}{ANSI_BLACK, "black foreground color"},
	// "\x1b[30mblack foreground color\x1b[0m"},
	// {[]interface{}{ANSI_BG_YELLOW, ANSI_BOLD, ANSI_RED,
	// "bold red text with yellow background",
	// ANSI_BG_CYAN, ANSI_GREEN, "green text with cyan background"},
	// "\x1b[43m\x1b[1m\x1b[31mbold red text with yellow background" +
	// "\x1b[46m\x1b[32mgreen text with cyan background\x1b[0m"},
	// {[]interface{}{ANSI_UNDERLINE, ANSI_GREEN, "green underline text"},
	// "\x1b[4m\x1b[32mgreen underline text\x1b[0m"},
	// {[]interface{}{ANSI_BOLD, ANSI_UNDERLINE, ANSI_GREEN,
	// "bold green underline text"},
	// "\x1b[1m\x1b[4m\x1b[32mbold green underline text\x1b[0m"},
	// {[]interface{}{ANSI_BOLD, ANSI_GREEN, "colored ", ANSI_OFF,
	// "to normal text"},
	// "\x1b[1m\x1b[32mcolored \x1b[0mto normal text\x1b[0m"},
	// {[]interface{}{ANSI_BOLD, ANSI_GREEN, "colored", ANSI_OFF},
	// "\x1b[1m\x1b[32mcolored\x1b[0m"},
}

func TestColors(t *testing.T) {
	for i, v := range colorTests {
		if out := v.escapeCodes; out != v.output {
			fmt.Println(v.escapeCodes)
			t.Errorf("Test Number: %d\nGot:\t%q\nExpect:\t%q\n", i,
				v.escapeCodes, v.output)
		}
	}
}
