============================================
go-elog - An enhanced logging library for Go
============================================

go-elog is a simple logging library for Go that aims to be a more robust
logging package for Go than the log package in the standard library.

----------------
go-elog Features
----------------

Logging levels
==============

Never again have to hunt for and remove printf statements (or erroneous package
imports) when debugging code. Simply switch to another logging level.
Supported log levels include DEBUG, INFO, WARNING, ERROR, and CRITICAL.

ANSI text attributes
====================

Colorize and embolden your logging output with AnsiEscape(). Build complex
colore text output simply and easily by using:

.. code-block:: go

    AnsiEscape(BOLD, GREEN, BG_RED, "Bold green text with a red background")

Multiple output streams
=======================

With go-elog it is possible to log to stdout, stderr, and a file at the same
time.

Customize the output format with templates
==========================================

Logging output is customizable with text templates.

-------
Example
-------

Here is a simple example for using go-elog to hide debug output when not
needed.

.. code-block:: go

    import "github.com/demizer/go-elog"

    log.Println("This message will be sent to stdout.")
    log.Debugln("This message will only be shown on stderr if the logging level is DEBUG!")

------------
Contributors
------------

See the AUTHORS file.
