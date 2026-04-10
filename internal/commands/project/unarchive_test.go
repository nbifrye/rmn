package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnarchiveCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/unarchive.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUnarchive(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"alpha"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Unarchived project alpha") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestUnarchiveCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUnarchive(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"alpha"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Status != "ok" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestUnarchiveCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdUnarchive(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"alpha"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
