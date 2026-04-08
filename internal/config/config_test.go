package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RedmineURL != "" || cfg.APIKey != "" {
		t.Errorf("expected empty config, got: %+v", cfg)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)

	data, _ := json.Marshal(Config{
		RedmineURL: "https://redmine.example.com",
		APIKey:     "secret-key",
	})
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RedmineURL != "https://redmine.example.com" {
		t.Errorf("expected URL, got: %s", cfg.RedmineURL)
	}
	if cfg.APIKey != "secret-key" {
		t.Errorf("expected API key, got: %s", cfg.APIKey)
	}
}

func TestLoad_CorruptJSON(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid"), 0o600)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := &Config{
		RedmineURL: "https://redmine.example.com",
		APIKey:     "my-key",
	}
	err := cfg.Save()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "rmn", "config.json")
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}

	// Verify content
	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.RedmineURL != cfg.RedmineURL || loaded.APIKey != cfg.APIKey {
		t.Errorf("loaded config doesn't match saved: %+v", loaded)
	}
}
