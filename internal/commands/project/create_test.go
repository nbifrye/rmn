package project

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
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Project api.ProjectCreateParams `json:"project"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.Project.Name != "New" || payload.Project.Identifier != "new" {
			t.Errorf("unexpected payload: %+v", payload)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 42, "name": "New", "identifier": "new"},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--name", "New", "--identifier", "new"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Created project #42") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 42, "name": "New", "identifier": "new"},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--name", "New", "--identifier", "new", "--public", "-d", "Desc", "--parent", "3"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p api.Project
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &p); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestCreateCommand_MissingName(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--identifier", "new"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCreateCommand_MissingIdentifier(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--name", "New"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing identifier")
	}
}

func TestCreateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"project": map[string]interface{}{"id": 1}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--name", "X", "--identifier", "x"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
