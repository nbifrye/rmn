package timeentry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteCommand_SuccessYesFlag(t *testing.T) {
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
	cmd.SetArgs([]string{"1", "-y"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Deleted time entry #1") {
		t.Errorf("expected success message, got: %s", out)
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
	cmd.SetArgs([]string{"1", "-y"})

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

func TestDeleteCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc", "-y"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid ID")
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
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Deleted time entry #1") {
		t.Errorf("expected deletion message, got: %s", out)
	}
}

func TestDeleteCommand_ConfirmNo(t *testing.T) {
	f := newNoServerFactory(t)
	f.IO.In = bytes.NewBufferString("n\n")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Cancelled") {
		t.Errorf("expected cancelled message, got: %s", out)
	}
}

func TestDeleteCommand_ConfirmReadError(t *testing.T) {
	f := newNoServerFactory(t)
	f.IO.In = bytes.NewBufferString("")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
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
	cmd.SetArgs([]string{"1", "-y"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
