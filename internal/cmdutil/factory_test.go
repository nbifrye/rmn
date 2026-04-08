package cmdutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFactory(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := NewFactory()
	if f.IO == nil {
		t.Fatal("expected IO to be set")
	}
	if f.Config == nil {
		t.Fatal("expected Config func to be set")
	}
	if f.APIClient == nil {
		t.Fatal("expected APIClient func to be set")
	}
}

func TestNewFactory_Config(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := NewFactory()
	cfg, err := f.Config()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RedmineURL != "" || cfg.APIKey != "" {
		t.Errorf("expected empty config, got: %+v", cfg)
	}
}

func TestNewFactory_APIClient_MissingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := NewFactory()
	_, err := f.APIClient()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
	if err.Error() != "not configured: run 'rmn auth login' to set up authentication" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewFactory_APIClient_WithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)
	data, _ := json.Marshal(map[string]string{
		"redmine_url": "https://redmine.example.com",
		"api_key":     "test-key",
	})
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)

	f := NewFactory()
	client, err := f.APIClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.BaseURL != "https://redmine.example.com" {
		t.Errorf("unexpected BaseURL: %s", client.BaseURL)
	}
}

func TestSetFlagOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	f := NewFactory()
	f.SetFlagOverrides("https://override.example.com", "override-key")

	if f.flagURL != "https://override.example.com" {
		t.Errorf("expected flagURL to be set, got: %s", f.flagURL)
	}
	if f.flagAPIKey != "override-key" {
		t.Errorf("expected flagAPIKey to be set, got: %s", f.flagAPIKey)
	}
}

func TestNewFactory_APIClient_ConfigLoadError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create config.json as a directory to trigger a read error
	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)
	os.MkdirAll(filepath.Join(dir, "config.json"), 0o755)

	f := NewFactory()
	_, err := f.APIClient()
	if err == nil {
		t.Fatal("expected error when config loading fails")
	}
}

func TestNewFactory_APIClient_FlagOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create config with one set of values
	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)
	data, _ := json.Marshal(map[string]string{
		"redmine_url": "https://original.example.com",
		"api_key":     "original-key",
	})
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)

	f := NewFactory()
	f.SetFlagOverrides("https://override.example.com", "override-key")

	client, err := f.APIClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL != "https://override.example.com" {
		t.Errorf("expected override URL, got: %s", client.BaseURL)
	}
	if client.APIKey != "override-key" {
		t.Errorf("expected override key, got: %s", client.APIKey)
	}
}
