package version

import (
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdVersion(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdVersion(f)
	if cmd.Use != "version" {
		t.Errorf("expected Use=version, got %s", cmd.Use)
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
