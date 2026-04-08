package issue

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteCommand_WithConfirmation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	// Simulate user typing "y\n" for confirmation
	f.IO.In = bytes.NewBufferString("y\n")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Deleted issue #42")) {
		t.Errorf("expected delete message, got: %s", out)
	}
}

func TestDeleteCommand_WithYesFlag(t *testing.T) {
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
	cmd.SetArgs([]string{"42", "--yes"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Deleted issue #42")) {
		t.Errorf("expected delete message, got: %s", out)
	}
}

func TestDeleteCommand_Cancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make API call when cancelled")
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	f.IO.In = bytes.NewBufferString("n\n")
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Cancelled")) {
		t.Errorf("expected cancel message, got: %s", out)
	}
}

func TestDeleteCommand_InvalidID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}
