package commands

import (
	"bytes"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func newTestFactory() *cmdutil.Factory {
	return &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: "https://example.com", APIKey: "test"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient("https://example.com", "test"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}
}

func TestNewCmdRoot(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")

	if cmd.Use != "rmn" {
		t.Errorf("expected Use 'rmn', got %q", cmd.Use)
	}

	// Verify subcommands exist
	expected := []string{"auth", "issue", "mcp", "completion"}
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

func TestNewCmdRoot_PersistentPreRunE(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	// Use completion subcommand with flags to trigger PersistentPreRunE
	cmd.SetArgs([]string{"--redmine-url", "https://override.com", "--api-key", "new-key", "completion", "bash"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCmdCompletion_Bash(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"completion", "bash"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCmdCompletion_Zsh(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"completion", "zsh"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCmdCompletion_Fish(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"completion", "fish"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCmdCompletion_Powershell(t *testing.T) {
	f := newTestFactory()
	cmd := NewCmdRoot(f, "test")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"completion", "powershell"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
