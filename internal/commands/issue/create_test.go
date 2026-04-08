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

func TestCreateCommand_NumericProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		issue := body["issue"].(map[string]interface{})
		// Numeric string should be parsed to int
		if issue["project_id"] != float64(42) {
			t.Errorf("expected numeric project_id 42, got %v (type %T)", issue["project_id"], issue["project_id"])
		}
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{Issue: api.Issue{ID: 1, Subject: "test"}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "42", "--subject", "test"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCommand_AllFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		issue := body["issue"].(map[string]interface{})
		if issue["description"] != "A description" {
			t.Errorf("expected description, got %v", issue["description"])
		}
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{Issue: api.Issue{ID: 1, Subject: "test"}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{
		"--project", "test",
		"--subject", "Full issue",
		"--description", "A description",
		"--tracker", "1",
		"--priority", "2",
		"--assignee", "3",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--subject", "test"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API client failure")
	}
}

func TestCreateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":["Validation failed"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--subject", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestCreateCommand_JSONMarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{Issue: api.Issue{ID: 1, Subject: "test"}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	origMarshal := marshalJSON
	marshalJSON = func(v interface{}, prefix, indent string) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}
	defer func() { marshalJSON = origMarshal }()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--project", "test", "--subject", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for marshal failure")
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
