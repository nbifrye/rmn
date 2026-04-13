package project

import (
	"math"
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

// TestMarshalJSONInternalError covers the json.MarshalIndent error branch inside marshalJSON.
func TestMarshalJSONInternalError(t *testing.T) {
	// math.Inf produces a float64 that json.Marshal cannot encode.
	_, err := marshalJSON(math.Inf(1))
	if err == nil {
		t.Fatal("expected error marshaling infinity")
	}
}
