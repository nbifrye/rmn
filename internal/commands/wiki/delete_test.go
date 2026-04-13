package wiki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteCommand_YesFlag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "-y", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Deleted wiki page") {
		t.Errorf("expected success message")
	}
}

func TestDeleteCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "-y", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestDeleteCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-y", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteCommand_ConfirmYes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	f.IO.In = bytes.NewBufferString("y\n")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Deleted wiki page") {
		t.Errorf("expected success message")
	}
}

func TestDeleteCommand_ConfirmNo(t *testing.T) {
	f := newNoServerFactory(t)
	f.IO.In = bytes.NewBufferString("n\n")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Cancelled") {
		t.Errorf("expected cancelled message")
	}
}

func TestDeleteCommand_ConfirmReadError(t *testing.T) {
	f := newNoServerFactory(t)
	f.IO.In = bytes.NewBufferString("")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected read error")
	}
}

func TestDeleteCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "-y", "Start"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
