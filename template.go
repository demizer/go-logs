// Copyright 2013 The go-logger Authors. All rights reserved.
// This code is MIT licensed. See the LICENSE file for more info.

package logger

import "text/template"

// funcMap contains the available functions to the log format template.
var (
	funcMap = template.FuncMap{"ansiEscape": AnsiEscape}
	logFmt  = "{{if .Date}}{{.Date}} {{end}}" +
		"{{if .Prefix}}{{.Prefix}} {{end}}" +
		"{{if .LogLabel}}{{.LogLabel}} {{end}}" +
		"{{if .File}}{{.File}}:" +
		"{{if .Line}}{{.Line}}: {{end}}{{end}}" +
		"{{if .Text}}{{.Text}}{{end}}"
)

// format is the possible values that can be used in a log output format
type format struct {
	Prefix   string
	LogLabel string
	Date     string
	File     string
	Line     int
	Text     string
}
