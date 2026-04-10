package timeentry

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
			"time_entries": []map[string]interface{}{
				{"id": 1, "hours": 2.5, "comments": "Work", "spent_on": "2024-01-15",
					"project": map[string]interface{}{"id": 1, "name": "Test"},
					"user":    map[string]interface{}{"id": 1, "name": "Admin"},
					"activity": map[string]interface{}{"id": 1, "name": "Development"},
					"issue":   map[string]interface{}{"id": 10, "name": ""}},
			},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "2.50") {
		t.Errorf("expected hours in output, got: %s", out)
	}
	if !strings.Contains(out, "#10") {
		t.Errorf("expected issue reference in output, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entries": []map[string]interface{}{{"id": 1, "hours": 1.0}},
			"total_count":  1,
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result struct {
		TimeEntries []api.TimeEntry `json:"time_entries"`
		TotalCount  int             `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"time_entries": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "No time entries found.") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestListCommand_WithFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("project_id") != "test" {
			t.Errorf("expected project_id=test, got %s", r.URL.Query().Get("project_id"))
		}
		if r.URL.Query().Get("from") != "2024-01-01" {
			t.Errorf("expected from=2024-01-01, got %s", r.URL.Query().Get("from"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"time_entries": []map[string]interface{}{}, "total_count": 0})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--from", "2024-01-01"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
}
