package issue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
)

func TestNewCmdIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdIssue(f)

	if cmd.Use != "issue" {
		t.Errorf("expected Use 'issue', got %q", cmd.Use)
	}

	// Verify all 6 subcommands exist
	expected := []string{"list", "view", "create", "update", "close", "delete"}
	for _, name := range expected {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}

func TestMarshalJSON_Error(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	// Test that marshalJSON error propagates in create command (JSON output path)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{Issue: api.Issue{ID: 1, Subject: "test"}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--project", "1", "--subject", "test"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
	if err.Error() != "marshal error" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMarshalJSON_ErrorInView(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{Issue: api.Issue{ID: 1, Subject: "test"}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestMarshalJSON_ErrorInList(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: []api.Issue{{ID: 1}}, TotalCount: 1}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestMarshalJSON_ErrorInUpdate(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1", "--status", "3"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestMarshalJSON_ErrorInClose(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdClose(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestMarshalJSON_ErrorInDelete(t *testing.T) {
	original := marshalJSON
	defer func() { marshalJSON = original }()

	marshalJSON = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1", "--yes"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

func TestListCommand_FlushError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: []api.Issue{{ID: 1, Subject: "test"}}, TotalCount: 1}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	// Replace Out with an errWriter that fails after initial writes
	f.IO.Out = &errWriter{failAfter: 0}
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from flush")
	}
}

// errWriter is an io.Writer that returns an error after failAfter bytes.
type errWriter struct {
	written  int
	failAfter int
}

func (w *errWriter) Write(p []byte) (int, error) {
	if w.written >= w.failAfter {
		return 0, fmt.Errorf("write error")
	}
	w.written += len(p)
	if w.written > w.failAfter {
		return len(p), fmt.Errorf("write error")
	}
	return len(p), nil
}

func TestMarshalJSON_Default(t *testing.T) {
	// Verify the default marshalJSON works correctly
	data, err := marshalJSON(map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(data, []byte(`"key": "value"`)) {
		t.Errorf("unexpected JSON output: %s", string(data))
	}
}
