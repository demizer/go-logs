============
Things To Do
============
:Modified: Fri Oct 03 20:49 2014

-----
Items
-----

* Mon Sep 29 21:18 2014: Add log.WithFlags() \*Logger

  WithFlags(flags int) \*logger

  Allows logging a single output with flags passed as an argument.

* Sun Jul 13 19:20 2014: Add support for custom log control

  * x Tue Aug 12 20:57 2014 Tue Aug 12 20:58 2014: Logger.Excludes(functionNames ...string)

  * x Mon Aug 18 22:21 2014 Tue Aug 12 20:58 2014: Add ExcludeByHeirarchyID(xi ...int)

    Test without Ltree.

  * x Tue Aug 19 20:26 2014 Tue Aug 12 22:54 2014: Add ExcludeByString(xr ...string)

  * x Tue Aug 19 23:52 2014 Tue Aug 12 22:56 2014: Add ExcludeByFunctionName(xf ...string)

    Test with and without Lfilename flags.

  * Tue Aug 12 20:58 2014: Add ExcludeByRegexp(xr ...*Regexp)

* Tue Aug 19 23:55 2014: Label changes

  * x Sat Sep 06 10:09 2014 Tue Aug 12 20:39 2014: Rename [CRIT] to [CRITICAL] and the like.

  * Tue Aug 12 22:31 2014: Add special template for tree output

    The special template should keep the width of the labels to make the output
    uniform.

  * Tue Aug 12 20:40 2014: Allow setting of custom labels.

* Tue Aug 12 20:56 2014: Add output Show* functions.

  * Tue Aug 12 23:00 2014: Add ShowByHeirarchyID(sh ...int)

  * Tue Aug 12 22:59 2014: Add ShowByString(ss ...string)

  * Tue Aug 12 22:59 2014: Add ShowByFunctionName(sf ...string)

  * Tue Aug 12 22:58 2014: Add ShowByRegexp(sr ...*Regexp)

* Tue Aug 19 20:29 2014: Test logger.Write()

  There is a section where writing to a file can be tested.

* Tue Jul 15 03:48 2014: Write documentation

* Thu Jul 17 20:35 2014: Lint the code

* Tue Aug 12 22:52 2014: Split out Fprintf()... it's huge!

* Fri Jul 18 11:53 2014: Reorganize test functions in logger_test.go

* Mon Aug 18 22:07 2014: Remove duplicate test data

  For test tables, they should only be defined in logger_test.go, and not
  logger_std_test.go.

* Fri Oct 03 19:46 2014: If Lcolor is specified, then no color is generated.

  Currently stripAnsiColor() is used to remove the ansi coloring.

---------
Completed
---------

* x Fri Oct 03 20:46 2014 Mon Sep 29 23:23 2014: Test LnoFileAnsi with two streams of output

* x Fri Oct 03 20:47 2014 Thu Oct 02 22:17 2014: Fix SetStreams()

  SetStreams(f, os.Stdout)

  does not output to stdout.

* x Mon Sep 29 23:13 2014 Mon Sep 29 21:17 2014: Fix LnoFileAnsi

  It was not stopping ansi output from being dumped into files, and there seem
  to be no tests for it.

* x Mon Aug 18 22:19 2014 Mon Aug 18 22:15 2014: Add std = New() to std tests

  Prevents std from contaminating the next tests.

* x Sun Aug 10 22:53 2014 Sun Jul 13 22:26 2014: Add badges to README.rst

  * x Sun Aug 10 22:53 2014 Tue Jul 15 03:54 2014: Add example images

* x Sun Aug 10 12:42 2014 Fri Jul 18 00:56 2014: Add Llabel

* x Sun Aug 10 12:24 2014 Fri Jul 18 11:09 2014: Remove LEVEL_ALL

  The Print, Println, and Printf functions should always show output without
  LEVEL_ALL.

* Sun Jul 13 22:18 2014: 100% Test Coverage

  * x Thu Jul 17 23:58 2014 Thu Jul 17 22:34 2014: Add test for Ltree and Lindent cancellation.

    Ltree should be used even if Lindent is used.

  * x Thu Jul 17 22:33 2014 Thu Jul 17 22:33 2014: Test for Lindent and LshowIndent

  * x Fri Jul 18 00:03 2014 Thu Jul 17 22:34 2014: Test SetIndent(-2) when there is no indent

  * x Fri Jul 18 00:51 2014 Tue Jul 15 03:52 2014: Add test for bad templates

    log.SetTemplate("{{if .Date}}{{.Date}} {{end}}" +
        "{{if .Prefix}}{{.Prefix}} {{end}}" +
        "{{if .LogLabel}}{{.LogLabel}} {{end}}" +
        "{{if .Ident}}{{.Ident}}{{end}}" +
        "{{if .Id}}{{.Id}}{{end}}" +
        "{{if .FileName}}{{.FileName}}: {{end}}" +
        "{{if .FunctionName}}{{.FunctionName}}{{end}}" +
        "{{if .LineNumber}}#{{.LineNumber}}: {{end}}" +
        "{{if .Text}}{{.Text}}{{end}}")

    .Ident is suppossed to be .Indent. The incorrect occurrenc causes all
    output to be broken.

  * x Sun Aug 10 12:23 2014 Thu Jul 17 22:34 2014: Test remaining functions

* x Thu Jul 17 20:11 2014 Wed Jul 16 03:06 2014: Add column stop for emit() output so it lines up.

* x Thu Jul 17 20:11 2014 Wed Jul 16 03:01 2014: Add Ldots for indentation showing

  * x Thu Jul 17 20:11 2014 Change dots color to a grayesh blue

* Sun Jul 13 19:02 2014: Add support for 256 Colors

    http://en.wikipedia.org/wiki/ANSI_escape_code
    http://ascii-table.com/ansi-escape-sequences.php

   * x Mon Jul 14 00:31 2014 Remove old color functions

   * x Mon Jul 14 00:31 2014 Add https://github.com/Knorkebrot/ansirgb or https://github.com/aybabtme/rgbterm

   * x Mon Jul 14 20:39 2014 Change Lansi to Lcolor

* x Sun Jul 13 22:18 2014 Sun Jul 13 20:26 2014: Test LStdFlags

  It is not producing colored output for some reason.
