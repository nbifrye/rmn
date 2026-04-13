package version

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

// errWriter is an io.Writer that returns an error after failAfter bytes have been written.
type errWriter struct {
	written   int
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

// ---------------------------------------------------------------------------
// version.go — marshalJSON real error path (unmarshalable type)
// ---------------------------------------------------------------------------

func TestMarshalJSONVar_RealError(t *testing.T) {
	_, err := marshalJSON(make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable value")
	}
}

// ---------------------------------------------------------------------------
// list.go — marshalJSON error in JSON output mode
// ---------------------------------------------------------------------------

func TestListCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions:   []api.Version{{ID: 1, Name: "v1.0", Status: "open", Sharing: "none"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"--project", "test"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error from marshalJSON")
	}
}

// ---------------------------------------------------------------------------
// list.go — flush error (tabwriter)
// ---------------------------------------------------------------------------

func TestListCommand_FlushError(t *testing.T) {
	dueDate := "2024-06-30"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{
			Versions:   []api.Version{{ID: 1, Name: "v1.0", Status: "open", DueDate: &dueDate, Sharing: "none"}},
			TotalCount: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	f.IO.Out = &errWriter{failAfter: 0}
	cmd := NewCmdList(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error from flush")
	}
}

// ---------------------------------------------------------------------------
// create.go — f.APIClient() error
// ---------------------------------------------------------------------------

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
	cmd.SetArgs([]string{"--project", "test", "--name", "v1.0"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for APIClient failure")
	}
}

// ---------------------------------------------------------------------------
// create.go — client.CreateVersion API error (server 500)
// ---------------------------------------------------------------------------

func TestCreateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--project", "test", "--name", "v1.0"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for API 500")
	}
}

// ---------------------------------------------------------------------------
// update.go — f.APIClient() error
// ---------------------------------------------------------------------------

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
	cmd.SetArgs([]string{"1", "--name", "v1.1"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for APIClient failure")
	}
}

// ---------------------------------------------------------------------------
// update.go — client.UpdateVersion API error (server 500)
// ---------------------------------------------------------------------------

func TestUpdateCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdUpdate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--name", "v1.1"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for API 500")
	}
}

// ---------------------------------------------------------------------------
// delete.go — f.APIClient() error (after --yes)
// ---------------------------------------------------------------------------

func TestDeleteCommand_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{In: &bytes.Buffer{}, Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
	}
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--yes"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for APIClient failure")
	}
}

// ---------------------------------------------------------------------------
// delete.go — client.DeleteVersion API error (server 500)
// ---------------------------------------------------------------------------

func TestDeleteCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdDelete(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1", "--yes"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for API 500")
	}
}

// ---------------------------------------------------------------------------
// view.go — client.GetVersion API error (server 500)
// ---------------------------------------------------------------------------

func TestViewCommand_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdView(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for API 500")
	}
}
