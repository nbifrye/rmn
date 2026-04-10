package user

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
		if r.URL.Path != "/users/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{
				ID:          42,
				Login:       "jdoe",
				FirstName:   "John",
				LastName:    "Doe",
				Mail:        "jdoe@example.com",
				Admin:       false,
				LastLoginOn: "2024-01-15T10:00:00Z",
				CreatedOn:   "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"User #42", "jdoe", "John Doe", "jdoe@example.com", "No", "2024-01-15T10:00:00Z", "2024-01-01T00:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{ID: 42, Login: "jdoe", FirstName: "John", LastName: "Doe"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result api.User
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.ID != 42 {
		t.Errorf("expected ID 42, got %d", result.ID)
	}
}

func TestViewCommand_Me(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/current.json" {
			t.Errorf("expected /users/current.json, got: %s", r.URL.Path)
		}
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{
				ID:        1,
				Login:     "me",
				FirstName: "Current",
				LastName:  "User",
				Mail:      "me@example.com",
				Admin:     true,
				CreatedOn: "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"me"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	for _, want := range []string{"User #1", "Current User", "me@example.com", "Yes"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

func TestViewCommand_MeJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/current.json" {
			t.Errorf("expected /users/current.json, got: %s", r.URL.Path)
		}
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{ID: 1, Login: "me", FirstName: "Current", LastName: "User"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"me"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result api.User
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestViewCommand_InvalidID(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"abc"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
	if err.Error() != "invalid user ID: abc" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestViewCommand_AdminYes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{ID: 1, Login: "admin", FirstName: "Admin", LastName: "User", Admin: true, CreatedOn: "2024-01-01T00:00:00Z"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !strings.Contains(out, "Admin:       Yes") {
		t.Errorf("expected 'Admin:       Yes' in output, got: %s", out)
	}
}

func TestViewCommand_NoLastLogin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			User api.User `json:"user"`
		}{
			User: api.User{ID: 1, Login: "new", FirstName: "New", LastName: "User", CreatedOn: "2024-01-01T00:00:00Z"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if strings.Contains(out, "Last Login:") {
		t.Errorf("expected no 'Last Login:' line when last_login_on is empty, got: %s", out)
	}
}

func TestViewCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"42"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for API client failure")
	}
}

func TestViewCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"999"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestViewCommand_MissingArgs(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}
