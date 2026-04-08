package auth

import (
	"bytes"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestNewCmdAuth(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
		APIClient: func() (*api.Client, error) {
			return nil, nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdAuth(f)

	if cmd.Use != "auth" {
		t.Errorf("expected Use 'auth', got %q", cmd.Use)
	}

	// Verify subcommands exist
	expected := []string{"login", "status"}
	for _, name := range expected {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}
