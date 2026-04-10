package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
)

func TestCreateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/projects/alpha/versions.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Version api.VersionCreateParams `json:"version"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.Version.Name != "v1.0" {
			t.Errorf("unexpected payload: %+v", payload)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": map[string]interface{}{"id": 42, "name": "v1.0"},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--name", "v1.0", "--status", "open", "--sharing", "none", "--due-date", "2024-06-30", "-d", "desc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Created version #42") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"version": map[string]interface{}{"id": 1, "name": "v1"}})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--name", "v1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var v api.Version
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestCreateCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--name", "v1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCreateCommand_MissingName(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCreateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"version": map[string]interface{}{"id": 1}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--name", "v1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
