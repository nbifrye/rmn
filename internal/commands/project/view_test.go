package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestViewCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/projects/alpha.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{
				"id":          1,
				"name":        "Alpha",
				"identifier":  "alpha",
				"status":      1,
				"is_public":   true,
				"parent":      map[string]interface{}{"id": 5, "name": "Root"},
				"homepage":    "https://example.com",
				"created_on":  "2024-01-01",
				"updated_on":  "2024-02-01",
				"description": "Alpha project",
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"alpha"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"Project #1", "Alpha", "alpha", "active", "yes", "Root", "https://example.com", "Alpha project"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 1, "name": "Alpha", "identifier": "alpha"},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p api.Project
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &p); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID=1, got %d", p.ID)
	}
}

func TestViewCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"project": map[string]interface{}{"id": 1}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestViewCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
