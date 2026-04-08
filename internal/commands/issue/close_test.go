package issue

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

func TestCloseCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"10"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API client failure")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected 'not configured' in error, got: %v", err)
	}
}

func TestCloseCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"999"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API failure")
	}
	if !strings.Contains(err.Error(), "Not found") {
		t.Errorf("expected 'Not found' in error, got: %v", err)
	}
}

func TestCloseCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"xyz"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}
