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

func TestUpdateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Membership api.MembershipUpdateParams `json:"membership"`
		}
		_ = json.Unmarshal(body, &payload)
		if len(payload.Membership.RoleIDs) != 2 {
			t.Errorf("unexpected payload: %+v", payload)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--role", "3,4"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Updated membership #1") {
		t.Errorf("expected success message")
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
	cmd.SetArgs([]string{"1", "--role", "3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc", "--role", "3"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCommand_MissingRole(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
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
	cmd.SetArgs([]string{"1", "--role", "3"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--role", "1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--role", "1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
