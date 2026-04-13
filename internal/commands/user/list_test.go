package user

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
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users: []api.User{
				{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User", Mail: "admin@example.com", Admin: true},
				{ID: 2, Login: "jdoe", FirstName: "John", LastName: "Doe", Mail: "jdoe@example.com", Admin: false},
			},
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(resp)
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
	if out == "" {
		t.Error("expected output, got empty string")
	}
	for _, want := range []string{"admin", "Admin User", "admin@example.com", "Yes", "jdoe", "John Doe", "jdoe@example.com", "No"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
	if !strings.Contains(out, "Showing 2 of 2 users") {
		t.Errorf("expected output to contain summary, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users:      []api.User{{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
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
		Users      []api.User `json:"users"`
		TotalCount int        `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON output, got: %s", out)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
}

func TestListCommand_AdminColumn(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users: []api.User{
				{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User", Admin: true},
				{ID: 2, Login: "regular", FirstName: "Regular", LastName: "User", Admin: false},
			},
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(resp)
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
	if !strings.Contains(out, "Yes") {
		t.Errorf("expected 'Yes' for admin user, got: %s", out)
	}
	if !strings.Contains(out, "No") {
		t.Errorf("expected 'No' for non-admin user, got: %s", out)
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

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API failure")
	}
	if !strings.Contains(err.Error(), "Server error") {
		t.Errorf("expected 'Server error' in error, got: %v", err)
	}
}

func TestListCommand_WithFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "2" {
			t.Errorf("expected status=2, got %s", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("name") != "john" {
			t.Errorf("expected name=john, got %s", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("offset") != "5" {
			t.Errorf("expected offset=5, got %s", r.URL.Query().Get("offset"))
		}
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{Users: []api.User{}, TotalCount: 0}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--status", "2", "--name", "john", "--limit", "10", "--offset", "5"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users:      []api.User{},
			TotalCount: 0,
		}
		json.NewEncoder(w).Encode(resp)
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
	if !strings.Contains(out, "No users found.") {
		t.Errorf("expected 'No users found.' message, got: %s", out)
	}
}

func TestListCommand_MarshalJSONError(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("marshal error") }

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users:      []api.User{{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestListCommand_FlushError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{
			Users:      []api.User{{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	f.IO.Out = &errWriter{failAfter: 0}
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from flush")
	}
}

type errWriter struct {
	written   int
	failAfter int
}

func (w *errWriter) Write(p []byte) (int, error) {
	if w.written >= w.failAfter {
		return 0, fmt.Errorf("write error")
	}
	w.written += len(p)
	if w.written > w.failAfter {
		return len(p), fmt.Errorf("write error")
	}
	return len(p), nil
}
