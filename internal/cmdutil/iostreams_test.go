package cmdutil

import (
	"os"
	"testing"
)

func TestDefaultIOStreams(t *testing.T) {
	s := DefaultIOStreams()
	if s.In != os.Stdin {
		t.Error("expected In to be os.Stdin")
	}
	if s.Out != os.Stdout {
		t.Error("expected Out to be os.Stdout")
	}
	if s.ErrOut != os.Stderr {
		t.Error("expected ErrOut to be os.Stderr")
	}
}
