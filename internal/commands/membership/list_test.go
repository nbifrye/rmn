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

func TestListCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/projects/alpha/memberships.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"memberships": []map[string]interface{}{
				{
					"id":      1,
					"project": map[string]interface{}{"id": 10, "name": "Alpha"},
					"user":    map[string]interface{}{"id": 2, "name": "Alice"},
					"roles":   []map[string]interface{}{{"id": 3, "name": "Developer"}},
				},
				{
					"id":      2,
					"project": map[string]interface{}{"id": 10, "name": "Alpha"},
					"group":   map[string]interface{}{"id": 5, "name": "Devs"},
					"roles":   []map[string]interface{}{{"id": 4, "name": "Reporter"}},
				},
			},
			"total_count": 2,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"Alice", "Developer", "group: Devs", "Reporter", "Showing 2 of 2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"memberships": []map[string]interface{}{{"id": 1, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}}},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Memberships []api.Membership `json:"memberships"`
		TotalCount  int              `json:"total_count"`
	}
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestListCommand_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"memberships": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "No memberships found") {
		t.Errorf("expected empty message")
	}
}

func TestListCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestListCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"memberships": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestListCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
