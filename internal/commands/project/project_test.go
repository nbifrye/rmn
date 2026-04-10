package project

import (
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdProject(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdProject(f)
	if cmd.Use != "project" {
		t.Errorf("expected Use=project, got %s", cmd.Use)
	}
	subs := cmd.Commands()
	names := map[string]bool{}
	for _, c := range subs {
		names[c.Name()] = true
	}
	for _, want := range []string{"list", "view", "create", "update", "archive", "unarchive", "delete"} {
		if !names[want] {
			t.Errorf("expected subcommand %q to be registered", want)
		}
	}
}
