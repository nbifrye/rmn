package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func TestStatusCommand_Configured(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/current.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := struct {
			User struct {
				ID    int    `json:"id"`
				Login string `json:"login"`
			} `json:"user"`
		}{
			User: struct {
				ID    int    `json:"id"`
				Login string `json:"login"`
			}{ID: 1, Login: "admin"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: srv.URL, APIKey: "abcdef1234"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient(srv.URL, "abcdef1234"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Logged in as: admin")) {
		t.Errorf("expected login info, got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("abcd***")) {
		t.Errorf("expected masked API key, got: %s", out)
	}
}

func TestStatusCommand_NotConfigured(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("Not configured")) {
		t.Errorf("expected not-configured message, got: %s", out)
	}
}

func TestStatusCommand_ShortAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			User struct {
				ID    int    `json:"id"`
				Login string `json:"login"`
			} `json:"user"`
		}{
			User: struct {
				ID    int    `json:"id"`
				Login string `json:"login"`
			}{ID: 1, Login: "user"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: srv.URL, APIKey: "ab"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient(srv.URL, "ab"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	// Short key should show just "***"
	if !bytes.Contains([]byte(out), []byte("API Key:     ***")) {
		t.Errorf("expected masked short key, got: %s", out)
	}
}

func TestStatusCommand_EmptyAPIKey(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: "https://example.com", APIKey: ""}, nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("(not set)")) {
		t.Errorf("expected '(not set)' for empty key, got: %s", out)
	}
}

func TestStatusCommand_ConfigError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return nil, fmt.Errorf("config error")
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for config failure")
	}
}

func TestStatusCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: "https://example.com", APIKey: "some-key-12345"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("api client creation failed")
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error (should log, not return): %v", err)
	}

	errOut := f.IO.ErrOut.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(errOut), []byte("Connection failed")) {
		t.Errorf("expected connection failure in stderr, got: %s", errOut)
	}
}

func TestStatusCommand_ConnectionFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":["Unauthorized"]}`))
	}))
	defer srv.Close()

	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: srv.URL, APIKey: "bad-key-12345"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient(srv.URL, "bad-key-12345"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdStatus(f)
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errOut := f.IO.ErrOut.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(errOut), []byte("Connection failed")) {
		t.Errorf("expected connection failure message, got: %s", errOut)
	}
}
