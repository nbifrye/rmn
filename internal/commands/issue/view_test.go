package issue

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

func TestViewCommand_TableOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{
				ID:      42,
				Subject: "Test issue",
				Project: api.IdName{ID: 1, Name: "My Project"},
				Tracker: api.IdName{ID: 1, Name: "Bug"},
				Status:  api.IdName{ID: 1, Name: "New"},
				Priority: api.IdName{ID: 2, Name: "Normal"},
				Author:  api.IdName{ID: 1, Name: "Admin"},
				AssignedTo: &api.IdName{ID: 2, Name: "Developer"},
				Description: "A test description",
				DoneRatio:   50,
				CreatedOn:   "2024-01-01T00:00:00Z",
				UpdatedOn:   "2024-01-02T00:00:00Z",
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
	for _, want := range []string{"Issue #42", "Test issue", "My Project", "Bug", "New", "Normal", "Admin", "Developer", "50%", "A test description"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

func TestViewCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{ID: 42, Subject: "Test"},
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
	var result api.Issue
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("expected valid JSON, got: %s", out)
	}
	if result.ID != 42 {
		t.Errorf("expected ID 42, got %d", result.ID)
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
	if err.Error() != "invalid issue ID: abc" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestViewCommand_NoDescription(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{ID: 1, Subject: "No desc", Description: ""},
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
	if !bytes.Contains([]byte(out), []byte("No desc")) {
		t.Errorf("expected subject in output, got: %s", out)
	}
	// When description is empty, the production code skips the description section
	// (which would be a blank line followed by description text).
	// Verify the output ends at the "Updated:" line without trailing description content.
	if bytes.Contains([]byte(out), []byte("Updated:")) {
		// Find everything after the "Updated:" line
		idx := bytes.Index([]byte(out), []byte("Updated:"))
		afterUpdated := string([]byte(out)[idx:])
		// After "Updated:     <value>\n" there should be nothing else
		lines := bytes.Split([]byte(afterUpdated), []byte("\n"))
		// Expect: ["Updated:     ", ""] (the trailing newline produces an empty final element)
		nonEmpty := 0
		for _, line := range lines {
			if len(bytes.TrimSpace(line)) > 0 {
				nonEmpty++
			}
		}
		if nonEmpty != 1 {
			t.Errorf("expected no description section after Updated line, got: %s", out)
		}
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
	// No args provided
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestViewCommand_Unassigned(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Issue api.Issue `json:"issue"`
		}{
			Issue: api.Issue{ID: 1, Subject: "No assignee"},
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
	if !bytes.Contains([]byte(out), []byte("(unassigned)")) {
		t.Errorf("expected '(unassigned)' in output, got: %s", out)
	}
}
