package wiki

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
)

func TestCreateCommand_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			WikiPage api.WikiPageCreateParams `json:"wiki_page"`
		}
		_ = json.Unmarshal(body, &payload)
		if payload.WikiPage.Text != "Hello" {
			t.Errorf("unexpected payload: %+v", payload)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "NewPage", "text": "Hello", "version": 1, "author": map[string]interface{}{"id": 1, "name": "A"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "--text", "Hello", "-c", "first", "NewPage"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(f.IO.Out.(*bytes.Buffer).String(), "Created wiki page") {
		t.Errorf("expected success message")
	}
}

func TestCreateCommand_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "NewPage", "version": 1, "author": map[string]interface{}{"id": 1, "name": "A"}},
		})
	}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--text", "Hi", "NewPage"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCommand_MissingProject(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"--text", "Hi", "NewPage"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_MissingText(t *testing.T) {
	f := newNoServerFactory(t)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "table")
	cmd.SetArgs([]string{"-p", "alpha", "NewPage"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateCommand_MarshalJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"wiki_page": map[string]interface{}{"title": "X", "author": map[string]interface{}{"id": 1, "name": "A"}}})
	}))
	defer srv.Close()

	orig := marshalJSON
	defer func() { marshalJSON = orig }()
	marshalJSON = func(v interface{}) ([]byte, error) { return nil, fmt.Errorf("boom") }

	f := newTestFactory(srv)
	cmd := NewCmdCreate(f)
	setupRootFlags(cmd, "json")
	cmd.SetArgs([]string{"-p", "alpha", "--text", "Hi", "NewPage"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error")
	}
}
