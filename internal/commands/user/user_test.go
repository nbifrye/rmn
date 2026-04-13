package user

import (
	"bytes"
	"testing"

	"github.com/nbifrye/rmn/internal/cmdutil"
)

func TestNewCmdUser(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdUser(f)
	if cmd.Use != "user" {
		t.Errorf("expected Use=user, got %s", cmd.Use)
	}
	subs := cmd.Commands()
	names := map[string]bool{}
	for _, c := range subs {
		names[c.Name()] = true
	}
	for _, want := range []string{"list", "view"} {
		if !names[want] {
			t.Errorf("expected subcommand %q to be registered", want)
		}
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
