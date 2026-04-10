package wiki

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
		if !strings.Contains(r.URL.Path, "/projects/alpha/wiki/Start.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{
				"title":      "Start",
				"text":       "Hello world",
				"version":    3,
				"parent":     map[string]interface{}{"id": 1, "name": "Root"},
				"author":     map[string]interface{}{"id": 2, "name": "Alice"},
				"comments":   "minor edit",
				"created_on": "2024-01-01",
				"updated_on": "2024-02-01",
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"Start", "Hello world", "Alice", "Root", "minor edit"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestViewCommand_WithVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("version") != "2" {
			t.Errorf("expected version=2, got %s", r.URL.Query().Get("version"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "Start", "version": 2, "author": map[string]interface{}{"id": 1, "name": "A"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--version", "2", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "Start", "version": 1, "author": map[string]interface{}{"id": 1, "name": "A"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p api.WikiPageDetail
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &p); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestViewCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestViewCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"wiki_page": map[string]interface{}{"title": "Start", "author": map[string]interface{}{"id": 1, "name": "A"}}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
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
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
