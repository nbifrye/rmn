package main

import (
	"context"
	"os"
	"testing"
)

func TestRun_Help(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := run(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error running with no args: %v", err)
	}
}

func TestRun_Version(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := run(context.Background(), "1.0.0-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMain_Success(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	exitCode := -1
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = os.Exit }()

	oldArgs := os.Args
	os.Args = []string{"rmn"}
	defer func() { os.Args = oldArgs }()

	main()

	if exitCode != -1 {
		t.Errorf("expected no exit, got code %d", exitCode)
	}
}

func TestMain_Error(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	exitCode := -1
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = os.Exit }()

	oldArgs := os.Args
	os.Args = []string{"rmn", "issue", "view"}
	defer func() { os.Args = oldArgs }()

	main()

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}
