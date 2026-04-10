package timeentry

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

func TestCreateCommand_WithIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			TimeEntry api.TimeEntryCreateParams `json:"time_entry"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.TimeEntry.IssueID != 10 || payload.TimeEntry.Hours != 2.5 {
			t.Errorf("unexpected payload: %+v", payload)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 42, "hours": 2.5},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--issue", "10", "--hours", "2.5", "--activity", "9", "--spent-on", "2024-01-15", "-c", "Work"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Created time entry #42") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestCreateCommand_WithProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 1, "hours": 1.0},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--hours", "1.0"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 1, "hours": 1.0},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--issue", "10", "--hours", "1.0"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e api.TimeEntry
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestCreateCommand_MissingHours(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--issue", "10"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing hours")
	}
}

func TestCreateCommand_MissingIssueAndProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--hours", "1.0"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing issue and project")
	}
}

func TestCreateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"time_entry": map[string]interface{}{"id": 1}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--issue", "1", "--hours", "1.0"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
