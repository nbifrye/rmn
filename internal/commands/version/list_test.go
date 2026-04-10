package version

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
	dueDate := "2024-06-30"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-proj/versions.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions: []api.Version{
				{ID: 1, Name: "v1.0", Status: "open", DueDate: &dueDate, Sharing: "none"},
				{ID: 2, Name: "v2.0", Status: "locked", Sharing: "descendants"},
			},
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test-proj"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if out == "" {
		t.Error("expected output, got empty string")
	}
	for _, want := range []string{"v1.0", "v2.0", "open", "locked", "2024-06-30", "none", "descendants"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
	if !strings.Contains(out, "Showing 2 of 2 versions") {
		t.Errorf("expected output to contain summary, got: %s", out)
	}
}

func TestListCommand_NilDueDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions: []api.Version{
				{ID: 1, Name: "v1.0", Status: "open", DueDate: nil, Sharing: "none"},
			},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test-proj"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "v1.0") {
		t.Errorf("expected output to contain 'v1.0', got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions:   []api.Version{{ID: 1, Name: "v1.0", Status: "open", Sharing: "none"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--project", "test-proj"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result struct {
		Versions   []api.Version `json:"versions"`
		TotalCount int           `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON output, got: %s", out)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions:   []api.Version{},
			TotalCount: 0,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test-proj"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "No versions found.") {
		t.Errorf("expected 'No versions found.' message, got: %s", out)
	}
}

func TestListCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing project")
	}
	if !strings.Contains(err.Error(), "--project is required") {
		t.Errorf("expected '--project is required' in error, got: %v", err)
	}
}

func TestListCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API client failure")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected 'not configured' in error, got: %v", err)
	}
}

func TestListCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API failure")
	}
	if !strings.Contains(err.Error(), "Server error") {
		t.Errorf("expected 'Server error' in error, got: %v", err)
	}
}
