package membership

import (
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdMembership(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdMembership(f)
	if cmd.Use != "membership" {
		t.Errorf("expected Use=membership, got %s", cmd.Use)
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
