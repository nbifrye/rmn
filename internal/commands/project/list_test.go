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

func TestListCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects": []map[string]interface{}{
				{"id": 1, "name": "Alpha", "identifier": "alpha", "status": 1, "is_public": true},
				{"id": 2, "name": "Beta", "identifier": "beta", "status": 5, "is_public": false},
			},
			"total_count": 2,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "alpha") {
		t.Errorf("expected project in output, got: %s", out)
	}
	if !strings.Contains(out, "active") || !strings.Contains(out, "archived") {
		t.Errorf("expected status strings, got: %s", out)
	}
	if !strings.Contains(out, "Showing 2 of 2") {
		t.Errorf("expected total line, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects":    []map[string]interface{}{{"id": 1, "name": "Alpha", "identifier": "alpha"}},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result struct {
		Projects   []api.Project `json:"projects"`
		TotalCount int           `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.TotalCount != 1 || len(result.Projects) != 1 {
		t.Errorf("unexpected JSON result: %+v", result)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "No projects found.") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestListCommand_StatusFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("status"); got != "archived" {
			t.Errorf("expected status=archived, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--status", "archived"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")

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
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestListCommand_StatusStringFallback(t *testing.T) {
	if got := projectStatusString(42); got != "42" {
		t.Errorf("expected fallback to number, got %s", got)
	}
}
