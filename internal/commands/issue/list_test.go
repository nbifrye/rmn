package issue

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func newTestFactory(srv *httptest.Server) *cmdutil.Factory {
	return &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: srv.URL, APIKey: "test"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient(srv.URL, "test"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}
}

func TestListCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues: []api.Issue{
				{ID: 1, Subject: "First issue", Tracker: api.IdName{Name: "Bug"}, Status: api.IdName{Name: "Open"}, Priority: api.IdName{Name: "Normal"}},
				{ID: 2, Subject: "Second issue", Tracker: api.IdName{Name: "Feature"}, Status: api.IdName{Name: "Closed"}, Priority: api.IdName{Name: "High"}},
			},
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	cmd.Root().PersistentFlags().String("output", "table", "")
	cmd.Root().PersistentFlags().String("redmine-url", "", "")
	cmd.Root().PersistentFlags().String("api-key", "", "")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if out == "" {
		t.Error("expected output, got empty string")
	}
	if !bytes.Contains([]byte(out), []byte("First issue")) {
		t.Errorf("expected output to contain 'First issue', got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("Showing 2 of 2 issues")) {
		t.Errorf("expected output to contain summary, got: %s", out)
	}
}

func TestListCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues:     []api.Issue{{ID: 1, Subject: "Test"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	cmd.Root().PersistentFlags().String("output", "json", "")
	cmd.Root().PersistentFlags().String("redmine-url", "", "")
	cmd.Root().PersistentFlags().String("api-key", "", "")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	var result struct {
		Issues     []api.Issue `json:"issues"`
		TotalCount int         `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON output, got: %s", out)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
}

func TestListCommand_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{
			Issues:     []api.Issue{},
			TotalCount: 0,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	cmd.Root().PersistentFlags().String("output", "table", "")
	cmd.Root().PersistentFlags().String("redmine-url", "", "")
	cmd.Root().PersistentFlags().String("api-key", "", "")

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := f.IO.Out.(*bytes.Buffer).String()
	if !bytes.Contains([]byte(out), []byte("No issues found.")) {
		t.Errorf("expected 'No issues found.' message, got: %s", out)
	}
}
