package project

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
			Project api.ProjectUpdateParams `json:"project"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.Project.Name == nil || *payload.Project.Name != "Updated" {
			t.Errorf("expected name=Updated in payload, got %+v", payload)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"alpha", "--name", "Updated", "--public", "--parent", "3", "-d", "New desc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Updated project alpha") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestUpdateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1", "--name", "X"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Status != "ok" || result.ID != "1" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestUpdateCommand_NoFields(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"alpha"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errOut := f.IO.ErrOut.(*bytes.Buffer).String()
	if !strings.Contains(errOut, "No fields specified") {
		t.Errorf("expected no fields message, got: %s", errOut)
	}
}

func TestUpdateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1", "--name", "X"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
