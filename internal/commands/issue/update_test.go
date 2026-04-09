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

func TestUpdateCommand_WithFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Issue api.IssueUpdateParams `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.Notes != "test note" {
			t.Errorf("expected notes 'test note', got %q", body.Issue.Notes)
		}
		if body.Issue.StatusID == nil || *body.Issue.StatusID != 3 {
			t.Errorf("expected status_id 3, got %v", body.Issue.StatusID)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42", "--status", "3", "--notes", "test note"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Updated issue #42")) {
		t.Errorf("expected update message, got: %s", out)
	}
}

func TestUpdateCommand_NoFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make API call when no flags are set")
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errOut := f.IO.ErrOut.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(errOut), []byte("No fields specified")) {
		t.Errorf("expected no-fields warning, got: %s", errOut)
	}
}

func TestUpdateCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc", "--status", "1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}

func TestUpdateCommand_AllFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Issue map[string]interface{} `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue["subject"] != "New subject" {
			t.Errorf("expected subject 'New subject', got %v", body.Issue["subject"])
		}
		if body.Issue["notes"] != "A note" {
			t.Errorf("expected notes 'A note', got %v", body.Issue["notes"])
		}
		if body.Issue["status_id"] != float64(3) {
			t.Errorf("expected status_id 3, got %v", body.Issue["status_id"])
		}
		if body.Issue["tracker_id"] != float64(2) {
			t.Errorf("expected tracker_id 2, got %v", body.Issue["tracker_id"])
		}
		if body.Issue["priority_id"] != float64(1) {
			t.Errorf("expected priority_id 1, got %v", body.Issue["priority_id"])
		}
		if body.Issue["description"] != "New description" {
			t.Errorf("expected description 'New description', got %v", body.Issue["description"])
		}
		if body.Issue["assigned_to_id"] != float64(5) {
			t.Errorf("expected assigned_to_id 5, got %v", body.Issue["assigned_to_id"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42",
		"--status", "3",
		"--tracker", "2",
		"--priority", "1",
		"--subject", "New subject",
		"--description", "New description",
		"--assignee", "5",
		"--notes", "A note",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42", "--status", "3"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API client failure")
	}
}

func TestUpdateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42", "--status", "3"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestUpdateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"42", "--status", "3"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result struct {
		Status  string `json:"status"`
		ID      int    `json:"id"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", result.Status)
	}
	if result.ID != 42 {
		t.Errorf("expected id 42, got %d", result.ID)
	}
}

func TestUpdateCommand_ExtendedFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Issue map[string]interface{} `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue["start_date"] != "2024-06-01" {
			t.Errorf("expected start_date '2024-06-01', got %v", body.Issue["start_date"])
		}
		if body.Issue["due_date"] != "2024-12-31" {
			t.Errorf("expected due_date '2024-12-31', got %v", body.Issue["due_date"])
		}
		if body.Issue["estimated_hours"] != float64(16) {
			t.Errorf("expected estimated_hours 16, got %v", body.Issue["estimated_hours"])
		}
		if body.Issue["done_ratio"] != float64(75) {
			t.Errorf("expected done_ratio 75, got %v", body.Issue["done_ratio"])
		}
		if body.Issue["category_id"] != float64(5) {
			t.Errorf("expected category_id 5, got %v", body.Issue["category_id"])
		}
		if body.Issue["fixed_version_id"] != float64(2) {
			t.Errorf("expected fixed_version_id 2, got %v", body.Issue["fixed_version_id"])
		}
		if body.Issue["parent_issue_id"] != float64(20) {
			t.Errorf("expected parent_issue_id 20, got %v", body.Issue["parent_issue_id"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42",
		"--start-date", "2024-06-01",
		"--due-date", "2024-12-31",
		"--estimated-hours", "16",
		"--done-ratio", "75",
		"--category", "5",
		"--version", "2",
		"--parent", "20",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCommand_SubjectOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Issue map[string]interface{} `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue["subject"] != "New subject" {
			t.Errorf("expected subject 'New subject', got %v", body.Issue["subject"])
		}
		// status_id should not be present
		if _, exists := body.Issue["status_id"]; exists {
			t.Error("status_id should not be sent when not specified")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42", "--subject", "New subject"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
