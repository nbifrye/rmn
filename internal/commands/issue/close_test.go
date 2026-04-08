package issue

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
)

func TestCloseCommand_DefaultStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/issues/10.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Issue api.IssueUpdateParams `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.StatusID == nil || *body.Issue.StatusID != defaultClosedStatusID {
			t.Errorf("expected status_id %d, got %v", defaultClosedStatusID, body.Issue.StatusID)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"10"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Closed issue #10")) {
		t.Errorf("expected close message, got: %s", out)
	}
}

func TestCloseCommand_CustomStatusAndNotes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Issue api.IssueUpdateParams `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.StatusID == nil || *body.Issue.StatusID != 6 {
			t.Errorf("expected status_id 6, got %v", body.Issue.StatusID)
		}
		if body.Issue.Notes != "Closing as won't fix" {
			t.Errorf("expected notes, got %q", body.Issue.Notes)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"10", "--status", "6", "--notes", "Closing as won't fix"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloseCommand_InvalidID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"xyz"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}
