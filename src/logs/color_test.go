// Copyright 2013,2015 The go-logs Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package logs

import (
	"fmt"
	"testing"

	"github.com/aybabtme/rgbterm"
)

var colorTests = []struct {
	escapeCodes string
	output      string
}{
	{rgbterm.FgString("red foreground color", 255, 0, 0),
		"\x1b[38;5;196mred foreground color\x1b[0;00m"},
	{rgbterm.FgString("green foreground color", 0, 255, 0),
		"\x1b[38;5;46mgreen foreground color\x1b[0;00m"},
	{rgbterm.FgString("blue foreground color", 0, 0, 255),
		"\x1b[38;5;21mblue foreground color\x1b[0;00m"},
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
