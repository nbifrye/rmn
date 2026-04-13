package membership

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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{
				"id":      5,
				"project": map[string]interface{}{"id": 10, "name": "Alpha"},
				"user":    map[string]interface{}{"id": 2, "name": "Alice"},
				"roles":   []map[string]interface{}{{"id": 3, "name": "Developer"}, {"id": 4, "name": "Manager"}},
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"5"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"Membership #5", "Alpha", "Alice", "Developer", "Manager"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestViewCommand_TableOutput_Group(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{
				"id":      6,
				"project": map[string]interface{}{"id": 10, "name": "Alpha"},
				"group":   map[string]interface{}{"id": 7, "name": "Devs"},
				"roles":   []map[string]interface{}{{"id": 3, "name": "Developer"}},
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"6"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Devs") {
		t.Errorf("expected group name in output")
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{"id": 1, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}},
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

	var m api.Membership
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestViewCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestViewCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"membership": map[string]interface{}{"id": 1, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}}})
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

func TestViewCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
