package timeentry

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

func TestViewCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{
				"id":         7,
				"project":    map[string]interface{}{"id": 1, "name": "Alpha"},
				"issue":      map[string]interface{}{"id": 99},
				"user":       map[string]interface{}{"id": 2, "name": "Alice"},
				"activity":   map[string]interface{}{"id": 3, "name": "Development"},
				"hours":      2.5,
				"comments":   "Working",
				"spent_on":   "2024-01-15",
				"created_on": "2024-01-15T10:00:00Z",
				"updated_on": "2024-01-15T11:00:00Z",
			},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"7"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"Time Entry #7", "Alpha", "#99", "Alice", "Development", "2.50", "2024-01-15", "Working"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 1, "hours": 1.0,
				"project":  map[string]interface{}{"id": 1, "name": "Alpha"},
				"user":     map[string]interface{}{"id": 1, "name": "Admin"},
				"activity": map[string]interface{}{"id": 1, "name": "Dev"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e api.TimeEntry
	if err := json.Unmarshal(f.IO.Out.(*bytes.Buffer).Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestViewCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid ID")
	}
}

func TestViewCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"time_entry": map[string]interface{}{"id": 1}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestViewCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
