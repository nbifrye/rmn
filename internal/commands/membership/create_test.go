package membership

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
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestCreateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Membership api.MembershipCreateParams `json:"membership"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.Membership.UserID != 2 || len(payload.Membership.RoleIDs) != 2 {
			t.Errorf("unexpected payload: %+v", payload)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{"id": 42, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--user", "2", "--role", "3,4"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Created membership #42") {
		t.Errorf("expected success message")
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{"id": 1, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--user", "2", "--role", "3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--user", "2", "--role", "3"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_MissingUser(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--role", "3"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_MissingRole(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--user", "2"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"membership": map[string]interface{}{"id": 1, "project": map[string]interface{}{"id": 10, "name": "Alpha"}, "roles": []map[string]interface{}{}}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--user", "2", "--role", "3"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--user", "1", "--role", "1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--user", "1", "--role", "1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
