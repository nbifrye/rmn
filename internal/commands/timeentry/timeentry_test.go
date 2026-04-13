package timeentry

import (
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdTimeEntry(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdTimeEntry(f)
	if cmd.Use != "time-entry" {
		t.Errorf("expected Use=time-entry, got %s", cmd.Use)
	}
	names := map[string]bool{}
	for _, c := range cmd.Commands() {
		names[c.Name()] = true
	}
	for _, want := range []string{"list", "view", "create", "update", "delete"} {
		if !names[want] {
			t.Errorf("expected subcommand %q to be registered", want)
		}
	}
}
