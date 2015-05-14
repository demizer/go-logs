=================================
go-logs - Enhanced logging for Go
=================================

.. image:: https://travis-ci.org/demizer/go-logs.png?branch=master
    :target: https://travis-ci.org/demizer/go-logs
.. image:: https://img.shields.io/github.com/demizer/go-logs/status.png
    :target: https://drone.io/github.com/demizer/go-logs/latest
.. image:: https://coveralls.io/repos/demizer/go-logs/badge.png?branch=master
    :target: https://coveralls.io/r/demizer/go-logs?branch=master
.. image:: https://godoc.org/github.com/demizer/go-logs?status.svg
    :target: http://godoc.org/github.com/demizer/go-logs
|

A drop-in replacement for the Go standard library logger package.

.. image:: screenshot-2015-05-13.png

--------
Features
--------

* Logging levels
* Colored text output
* Multiple simultaneous output streams
* Customizable output formatting using templates
* Hierarchical output formatting
* Suppress specific output

-------
Example
-------

Here is a simple example for using go-logs to hide debug output when not
needed.

.. code-block:: go

    import "github.com/demizer/go-logs/src/logs"

    log.Println("This message will be sent to stdout.")
    log.Debugln("This message will only be shown on stderr if the logging level is DEBUG!")
