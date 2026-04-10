package wiki

import (
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdWiki(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdWiki(f)
	if cmd.Use != "wiki" {
		t.Errorf("expected Use=wiki, got %s", cmd.Use)
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
