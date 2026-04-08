package issue

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
)

func TestCreateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body struct {
			Issue api.IssueCreateParams `json:"issue"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.Subject != "New bug" {
			t.Errorf("unexpected subject: %v", body.Issue.Subject)
		}
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{ID: 99, Subject: "New bug"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test-proj", "--subject", "New bug"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Created issue #99")) {
		t.Errorf("expected creation message, got: %s", out)
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{ID: 99, Subject: "New bug"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--project", "1", "--subject", "New bug"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result api.Issue
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.ID != 99 {
		t.Errorf("expected ID 99, got %d", result.ID)
	}
}

func TestCreateCommand_MissingProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--subject", "No project"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCreateCommand_MissingSubject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing subject")
	}
}
