package wiki

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

func TestUpdateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			WikiPage api.WikiPageUpdateParams `json:"wiki_page"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.WikiPage.Text == nil || *payload.WikiPage.Text != "Updated" {
			t.Errorf("expected text=Updated, got %+v", payload)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--text", "Updated", "-c", "edit", "--version", "3", "Start"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Updated wiki page") {
		t.Errorf("expected success message")
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
	cmd.SetArgs([]string{"-p", "alpha", "--text", "New", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Status string `json:"status"`
		Title  string `json:"title"`
	}
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Title != "Start" {
		t.Errorf("unexpected title: %s", result.Title)
	}
}

func TestUpdateCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--text", "X", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCommand_NoFields(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(f.IO.ErrOut.(*bytes.Buffer).String(), "No fields specified") {
		t.Errorf("expected message")
	}
}

func TestUpdateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--text", "New", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
