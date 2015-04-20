package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/demizer/go-elog"
)

func main() {
	var buf bytes.Buffer

	logr := log.New(log.LEVEL_DEBUG, os.Stdout)
	logr.SetFlags(log.Ldate | log.Lseperator)

	logr.Println("\nDUAL STREAM OUTPUT EXAMPLE (like the tee command)")

	logr.SetFlags(log.LdebugFlags | log.Ldate | log.Lseperator)

	logr.Println("\nstdout output:\n")

	logr.SetStreams(os.Stdout, &buf)

	logr.Debugln("Level 0 Output 1")
	lvl3 := func() {
		logr.Debugln("Level 3 Output 1")
	}
	lvl2 := func() {
		logr.Debugln("Level 2 Output 1")
		logr.Criticalln("Level 2 Output 2")
		lvl3()
		logr.Debugln("Level 2 Output 3")
	}
	lvl1 := func() {
		logr.Infoln("Level 1 Output 1")
		logr.Errorln("Level 1 Output 2")
		lvl2()
		logr.Warningln("Level 1 Output 3")
	}
	lvl1()

	logr.Debugln("Level 0 Output 2")

	logr.SetStreams(os.Stdout)

	logr.Println("\nShowing output stored in the buffer:\n")

	fmt.Print(buf.String())
}
