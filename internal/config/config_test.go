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

func TestLoad_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	dir := filepath.Join(tmpDir, "rmn")
	os.MkdirAll(dir, 0o755)
	// Create config.json as a directory — ReadFile will fail with a non-ErrNotExist error
	os.MkdirAll(filepath.Join(dir, "config.json"), 0o755)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for unreadable file")
	}
}

func TestConfigDir_FallbackToHome(t *testing.T) {
	// Unset XDG_CONFIG_HOME to test fallback to HOME
	t.Setenv("XDG_CONFIG_HOME", "")

	dir, err := configDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Fatal("expected non-empty config dir")
	}
	// Should end with .config/rmn
	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got: %s", dir)
	}
}

func TestConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	path, err := configPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(tmpDir, "rmn", "config.json")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestConfigDir_HomeDirError(t *testing.T) {
	// Unset both XDG_CONFIG_HOME and HOME to trigger UserHomeDir error
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	_, err := configDir()
	if err == nil {
		t.Fatal("expected error when HOME is unset")
	}
}

func TestConfigPath_Error(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	_, err := configPath()
	if err == nil {
		t.Fatal("expected error when HOME is unset")
	}
}

func TestLoad_ConfigPathError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when HOME is unset")
	}
}

func TestSave_ConfigPathError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	cfg := &Config{RedmineURL: "https://example.com", APIKey: "key"}
	err := cfg.Save()
	if err == nil {
		t.Fatal("expected error when HOME is unset")
	}
}

func TestSave_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create a file where the directory should be, so MkdirAll fails
	os.WriteFile(filepath.Join(tmpDir, "rmn"), []byte("not a dir"), 0o644)

	cfg := &Config{RedmineURL: "https://example.com", APIKey: "key"}
	err := cfg.Save()
	if err == nil {
		t.Fatal("expected error when dir creation fails")
	}
}

func TestSave_OverwriteExisting(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := &Config{RedmineURL: "https://first.com", APIKey: "key1"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg2 := &Config{RedmineURL: "https://second.com", APIKey: "key2"}
	if err := cfg2.Save(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.RedmineURL != "https://second.com" {
		t.Errorf("expected overwritten URL, got: %s", loaded.RedmineURL)
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
