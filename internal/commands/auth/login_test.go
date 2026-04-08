package auth

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestLoginCommand_WithFlags(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) { return nil, nil },
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdLogin(f)
	cmd.SetArgs([]string{"--url", "https://redmine.example.com", "--api-key", "test-key-123"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Authentication configured successfully")) {
		t.Errorf("expected success message, got: %s", out)
	}

	// Verify config file was created
	cfgPath := filepath.Join(tmpDir, "rmn", "config.json")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

func TestLoginCommand_WithPrompts(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) { return nil, nil },
		IO: &cmdutil.IOStreams{
			In:     bytes.NewBufferString("https://redmine.example.com\nmy-api-key\n"),
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdLogin(f)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Authentication configured successfully")) {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestLoginCommand_EmptyInput(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) { return nil, nil },
		IO: &cmdutil.IOStreams{
			In:     bytes.NewBufferString("\n\n"),
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdLogin(f)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if err.Error() != "both Redmine URL and API key are required" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoginCommand_TrimsTrailingSlash(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) { return nil, nil },
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdLogin(f)
	cmd.SetArgs([]string{"--url", "https://redmine.example.com/", "--api-key", "key"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Load and verify the URL was trimmed
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.RedmineURL != "https://redmine.example.com" {
		t.Errorf("expected trailing slash trimmed, got: %s", cfg.RedmineURL)
	}
}
