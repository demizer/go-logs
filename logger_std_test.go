// Copyright 2013,2014 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

// Tests for the default standard logging object

package log

import (
	"bytes"
	"testing"
	"time"
)

func TestStdTemplate(t *testing.T) {
	var buf bytes.Buffer

	SetLevel(LEVEL_DEBUG)
	SetStreams(&buf)

	SetFlags(LdebugFlags)

	SetTemplate("{{.Text}}")
	temp := Template()

	type test struct {
		Text string
	}

	err := temp.Execute(&buf, &test{"Hello, World!"})
	if err != nil {
		t.Fatal(err)
	}

	expe := "Hello, World!"

	if buf.String() != expe {
		t.Errorf("\nGot:\t%s\nExpect:\t%s\n", buf.String(), expe)
	}
}

func TestStdSetTemplate(t *testing.T) {
	var buf bytes.Buffer

	SetLevel(LEVEL_DEBUG)
	SetStreams(&buf)

	SetFlags(LdebugFlags)

	SetTemplate("{{.Text}}")

	Debugln("Hello, World!")

	expe := "Hello, World!\n"

	if buf.String() != expe {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expe)
	}
}

func TestStdSetTemplateBad(t *testing.T) {
	var buf bytes.Buffer

	SetLevel(LEVEL_DEBUG)
	SetStreams(&buf)

	SetFlags(LdebugFlags)

	err := SetTemplate("{{.Text")

	Debugln("template: default:1: unclosed action")

	expe := "template: default:1: unclosed action"

	if err.Error() != expe {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expe)
	}
}

func TestStdSetTemplateBadDataObjectPanic(t *testing.T) {
	var buf bytes.Buffer

	SetLevel(LEVEL_DEBUG)

	SetStreams(&buf)

	SetFlags(LnoPrefix | Lindent)

	SetIndent(1)

	type test struct {
		Test string
	}

	err := SetTemplate("{{.Tes}}")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("\nGot:\t%q\nExpect:\tPANIC\n", buf.String())
		}
	}()

	Debugln("Hello, World!")

	// Reset the standard logging object
	SetTemplate(logFmt)
	SetIndent(0)
}

func TestStdDateFormat(t *testing.T) {
	dateFormat := DateFormat()

	expect := "Mon-20060102-15:04:05"

	if dateFormat != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", dateFormat, expect)
	}
}

func TestStdSetDateFormat(t *testing.T) {
	var buf bytes.Buffer

	SetLevel(LEVEL_ALL)

	SetStreams(&buf)

	SetFlags(Ldate)

	SetDateFormat("20060102-15:04:05")

	SetTemplate("{{.Date}}")

	Debugln("Hello")

	expect := time.Now().Format(DateFormat())

	if buf.String() != expect {
		t.Errorf("\nGot:\t%q\nExpect:\t%q\n", buf.String(), expect)
	}

	// Reset the standard logging object
	SetTemplate(logFmt)
}
