package tracker

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

func TestListCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"trackers": []map[string]interface{}{
				{"id": 1, "name": "Bug"},
				{"id": 2, "name": "Feature"},
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Bug") {
		t.Errorf("expected 'Bug' in output, got: %s", out)
	}
	if !strings.Contains(out, "Feature") {
		t.Errorf("expected 'Feature' in output, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"trackers": []map[string]interface{}{{"id": 1, "name": "Bug"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result []api.IdName
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"trackers": []map[string]interface{}{}})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "No trackers found.") {
		t.Errorf("expected 'No trackers found.' message, got: %s", out)
	}
}

func TestListCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}

	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
}
