package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	level := Level()
	if level != WARNING {
		t.Errorf("Level() = %s, want %s", level, WARNING)
	}
	SetLevel(DEBUG)
	level = Level()
	if level != DEBUG {
		t.Errorf("Level() = %s, want %s", level, DEBUG)
	}
}

func TestStream(t *testing.T) {
	if out := Stream(); out != os.Stderr {
		t.Errorf("Output stream is not stderr by default")
	}
	var buf bytes.Buffer
	SetStream(&buf)
	if out := Stream(); out != &buf {
		t.Errorf("SetOutput(&buf) = %p, want %p", out, &buf)
	}
}
