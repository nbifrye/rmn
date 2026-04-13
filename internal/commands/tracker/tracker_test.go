package tracker

import (
	"bytes"
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdTracker(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdTracker(f)
	if cmd.Use != "tracker" {
		t.Errorf("expected Use=tracker, got %s", cmd.Use)
	}
	subs := cmd.Commands()
	names := map[string]bool{}
	for _, c := range subs {
		names[c.Name()] = true
	}
	if !names["list"] {
		t.Error("expected subcommand 'list' to be registered")
	}
}

func TestMarshalJSON_Default(t *testing.T) {
	data, err := marshalJSON(map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(data, []byte(`"key": "value"`)) {
		t.Errorf("unexpected JSON output: %s", string(data))
	}
}

func TestMarshalJSON_RealError(t *testing.T) {
	_, err := marshalJSON(make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable value")
	}
}
