package issue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestListCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues: []api.Issue{
				{ID: 1, Subject: "First issue", Tracker: api.IdName{Name: "Bug"}, Status: api.IdName{Name: "Open"}, Priority: api.IdName{Name: "Normal"}},
				{ID: 2, Subject: "Second issue", Tracker: api.IdName{Name: "Feature"}, Status: api.IdName{Name: "Closed"}, Priority: api.IdName{Name: "High"}},
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
	if !bytes.Contains([]byte(out), []byte("First issue")) {
		t.Errorf("expected output to contain 'First issue', got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("Showing 2 of 2 issues")) {
		t.Errorf("expected output to contain summary, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues:     []api.Issue{{ID: 1, Subject: "Test"}},
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
		Issues     []api.Issue `json:"issues"`
		TotalCount int         `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON output, got: %s", out)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
}

func TestListCommand_WithAssignee(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues: []api.Issue{
				{
					ID:         1,
					Subject:    "Assigned issue",
					Tracker:    api.IdName{Name: "Bug"},
					Status:     api.IdName{Name: "Open"},
					Priority:   api.IdName{Name: "Normal"},
					AssignedTo: &api.IdName{ID: 2, Name: "Developer"},
				},
			},
			TotalCount: 1,
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
	if !bytes.Contains([]byte(out), []byte("Developer")) {
		t.Errorf("expected assignee name in output, got: %s", out)
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
}

func TestListCommand_WithFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("project_id") != "test" {
			t.Errorf("expected project_id=test, got %s", r.URL.Query().Get("project_id"))
		}
		if r.URL.Query().Get("status_id") != "closed" {
			t.Errorf("expected status_id=closed, got %s", r.URL.Query().Get("status_id"))
		}
		if r.URL.Query().Get("assigned_to_id") != "me" {
			t.Errorf("expected assigned_to_id=me, got %s", r.URL.Query().Get("assigned_to_id"))
		}
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: []api.Issue{}, TotalCount: 0}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--status", "closed", "--assignee", "me", "--tracker", "1", "--limit", "10", "--offset", "5"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListCommand_JSONMarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: []api.Issue{{ID: 1}}, TotalCount: 1}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	origMarshal := marshalJSON
	marshalJSON = func(v interface{}, prefix, indent string) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}
	defer func() { marshalJSON = origMarshal }()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for marshal failure")
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues:     []api.Issue{},
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
	if !bytes.Contains([]byte(out), []byte("No issues found.")) {
		t.Errorf("expected 'No issues found.' message, got: %s", out)
	}
}
