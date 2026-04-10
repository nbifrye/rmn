package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

// registerAndCapture creates a test MCP server, runs the given register
// function, and returns a map of registered tool handlers keyed by tool name.
func registerAndCapture(t *testing.T, register func(*mcp.Server, *api.Client), client *api.Client) map[string]mcp.ToolHandler {
	t.Helper()
	captured := make(map[string]mcp.ToolHandler)
	prev := capturedHandlers
	capturedHandlers = captured
	t.Cleanup(func() { capturedHandlers = prev })

	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.0"}, nil)
	register(s, client)
	return captured
}

func invokeTool(t *testing.T, handler mcp.ToolHandler, args map[string]interface{}) (string, bool) {
	t.Helper()
	argsJSON, _ := json.Marshal(args)
	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{Arguments: argsJSON},
	}
	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(result.Content) == 0 {
		t.Fatal("empty result content")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	return tc.Text, result.IsError
}

// --- Trackers & Statuses ---------------------------------------------------

func TestTrackerTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/trackers.json" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"trackers": []map[string]interface{}{{"id": 1, "name": "Bug"}},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")

	handlers := registerAndCapture(t, registerTrackerTools, client)
	text, isErr := invokeTool(t, handlers["list_trackers"], nil)
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
	if !strings.Contains(text, "Bug") {
		t.Errorf("expected tracker name in output: %s", text)
	}
}

func TestStatusTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/issue_statuses.json" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"issue_statuses": []map[string]interface{}{{"id": 1, "name": "New", "is_closed": false}},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")

	handlers := registerAndCapture(t, registerStatusTools, client)
	text, isErr := invokeTool(t, handlers["list_issue_statuses"], nil)
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
	if !strings.Contains(text, "New") {
		t.Errorf("expected status name in output: %s", text)
	}
}

// --- Projects --------------------------------------------------------------

func projectServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/projects.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"projects":    []map[string]interface{}{{"id": 1, "name": "Alpha", "identifier": "alpha"}},
				"total_count": 1,
			})
		case r.URL.Path == "/projects/alpha.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"project": map[string]interface{}{"id": 1, "name": "Alpha", "identifier": "alpha"},
			})
		case r.URL.Path == "/projects.json" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"project": map[string]interface{}{"id": 2, "name": "Beta", "identifier": "beta"},
			})
		case r.URL.Path == "/projects/alpha.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/projects/alpha/archive.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/projects/alpha/unarchive.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/projects/alpha.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestProjectTools(t *testing.T) {
	srv := projectServer()
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerProjectTools, client)

	if _, isErr := invokeTool(t, h["list_projects"], nil); isErr {
		t.Error("list_projects failed")
	}
	if _, isErr := invokeTool(t, h["get_project"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("get_project failed")
	}
	if _, isErr := invokeTool(t, h["get_project"], map[string]interface{}{}); !isErr {
		t.Error("expected error when project_id missing")
	}
	if _, isErr := invokeTool(t, h["create_project"], map[string]interface{}{"name": "Beta", "identifier": "beta"}); isErr {
		t.Error("create_project failed")
	}
	if _, isErr := invokeTool(t, h["create_project"], map[string]interface{}{"identifier": "beta"}); !isErr {
		t.Error("expected error when name missing")
	}
	if _, isErr := invokeTool(t, h["create_project"], map[string]interface{}{"name": "Beta"}); !isErr {
		t.Error("expected error when identifier missing")
	}
	if _, isErr := invokeTool(t, h["update_project"], map[string]interface{}{"project_id": "alpha", "name": "Alpha2", "is_public": true}); isErr {
		t.Error("update_project failed")
	}
	if _, isErr := invokeTool(t, h["update_project"], map[string]interface{}{}); !isErr {
		t.Error("expected error when project_id missing")
	}
	if _, isErr := invokeTool(t, h["archive_project"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("archive_project failed")
	}
	if _, isErr := invokeTool(t, h["archive_project"], map[string]interface{}{}); !isErr {
		t.Error("expected error when project_id missing")
	}
	if _, isErr := invokeTool(t, h["unarchive_project"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("unarchive_project failed")
	}
	if _, isErr := invokeTool(t, h["unarchive_project"], map[string]interface{}{}); !isErr {
		t.Error("expected error when project_id missing")
	}
	if _, isErr := invokeTool(t, h["delete_project"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("delete_project failed")
	}
	if _, isErr := invokeTool(t, h["delete_project"], map[string]interface{}{}); !isErr {
		t.Error("expected error when project_id missing")
	}
}

// --- Users -----------------------------------------------------------------

func TestUserTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users.json":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"users":       []map[string]interface{}{{"id": 1, "login": "alice"}},
				"total_count": 1,
			})
		case "/users/1.json":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user": map[string]interface{}{"id": 1, "login": "alice"},
			})
		case "/users/current.json":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user": map[string]interface{}{"id": 2, "login": "me"},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerUserTools, client)

	if _, isErr := invokeTool(t, h["list_users"], map[string]interface{}{"name": "ali"}); isErr {
		t.Error("list_users failed")
	}
	if _, isErr := invokeTool(t, h["get_user"], map[string]interface{}{"user_id": 1}); isErr {
		t.Error("get_user failed")
	}
	if _, isErr := invokeTool(t, h["get_user"], map[string]interface{}{"user_id": 0}); !isErr {
		t.Error("expected error for user_id 0")
	}
	if _, isErr := invokeTool(t, h["get_current_user"], nil); isErr {
		t.Error("get_current_user failed")
	}
}

// --- Versions --------------------------------------------------------------

func TestVersionTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/projects/alpha/versions.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"versions":    []map[string]interface{}{{"id": 10, "name": "v1"}},
				"total_count": 1,
			})
		case r.URL.Path == "/versions/10.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"version": map[string]interface{}{"id": 10, "name": "v1"},
			})
		case r.URL.Path == "/projects/alpha/versions.json" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"version": map[string]interface{}{"id": 11, "name": "v2"},
			})
		case r.URL.Path == "/versions/10.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/versions/10.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerVersionTools, client)

	if _, isErr := invokeTool(t, h["list_versions"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("list_versions failed")
	}
	if _, isErr := invokeTool(t, h["list_versions"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["get_version"], map[string]interface{}{"version_id": 10}); isErr {
		t.Error("get_version failed")
	}
	if _, isErr := invokeTool(t, h["get_version"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing version_id")
	}
	if _, isErr := invokeTool(t, h["create_version"], map[string]interface{}{"project_id": "alpha", "name": "v2"}); isErr {
		t.Error("create_version failed")
	}
	if _, isErr := invokeTool(t, h["create_version"], map[string]interface{}{"project_id": "alpha"}); !isErr {
		t.Error("expected error for missing name")
	}
	if _, isErr := invokeTool(t, h["create_version"], map[string]interface{}{"name": "v2"}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["update_version"], map[string]interface{}{"version_id": 10, "name": "v1.1"}); isErr {
		t.Error("update_version failed")
	}
	if _, isErr := invokeTool(t, h["update_version"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing version_id")
	}
	if _, isErr := invokeTool(t, h["delete_version"], map[string]interface{}{"version_id": 10}); isErr {
		t.Error("delete_version failed")
	}
	if _, isErr := invokeTool(t, h["delete_version"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing version_id")
	}
}

// --- Time Entries ----------------------------------------------------------

func TestTimeEntryTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/time_entries.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"time_entries": []map[string]interface{}{{"id": 5, "hours": 1.5}},
				"total_count":  1,
			})
		case r.URL.Path == "/time_entries/5.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"time_entry": map[string]interface{}{"id": 5, "hours": 1.5},
			})
		case r.URL.Path == "/time_entries.json" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"time_entry": map[string]interface{}{"id": 6, "hours": 2.0},
			})
		case r.URL.Path == "/time_entries/5.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/time_entries/5.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerTimeEntryTools, client)

	if _, isErr := invokeTool(t, h["list_time_entries"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("list_time_entries failed")
	}
	if _, isErr := invokeTool(t, h["get_time_entry"], map[string]interface{}{"time_entry_id": 5}); isErr {
		t.Error("get_time_entry failed")
	}
	if _, isErr := invokeTool(t, h["get_time_entry"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing time_entry_id")
	}
	if _, isErr := invokeTool(t, h["create_time_entry"], map[string]interface{}{"issue_id": 1, "hours": 2.0}); isErr {
		t.Error("create_time_entry failed")
	}
	if _, isErr := invokeTool(t, h["create_time_entry"], map[string]interface{}{"issue_id": 1}); !isErr {
		t.Error("expected error for missing hours")
	}
	if _, isErr := invokeTool(t, h["create_time_entry"], map[string]interface{}{"hours": 2.0}); !isErr {
		t.Error("expected error for missing issue_id/project_id")
	}
	if _, isErr := invokeTool(t, h["update_time_entry"], map[string]interface{}{"time_entry_id": 5, "hours": 3.0}); isErr {
		t.Error("update_time_entry failed")
	}
	if _, isErr := invokeTool(t, h["update_time_entry"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing time_entry_id")
	}
	if _, isErr := invokeTool(t, h["delete_time_entry"], map[string]interface{}{"time_entry_id": 5}); isErr {
		t.Error("delete_time_entry failed")
	}
	if _, isErr := invokeTool(t, h["delete_time_entry"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing time_entry_id")
	}
}

// --- Memberships -----------------------------------------------------------

func TestMembershipTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/projects/alpha/memberships.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"memberships": []map[string]interface{}{
					{"id": 1, "project": map[string]interface{}{"id": 1, "name": "Alpha"}, "user": map[string]interface{}{"id": 1, "name": "Alice"}, "roles": []map[string]interface{}{{"id": 3, "name": "Manager"}}},
				},
				"total_count": 1,
			})
		case r.URL.Path == "/memberships/1.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"membership": map[string]interface{}{"id": 1, "project": map[string]interface{}{"id": 1, "name": "Alpha"}, "user": map[string]interface{}{"id": 1, "name": "Alice"}},
			})
		case r.URL.Path == "/projects/alpha/memberships.json" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"membership": map[string]interface{}{"id": 2},
			})
		case r.URL.Path == "/memberships/1.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/memberships/1.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerMembershipTools, client)

	if _, isErr := invokeTool(t, h["list_memberships"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("list_memberships failed")
	}
	if _, isErr := invokeTool(t, h["list_memberships"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["get_membership"], map[string]interface{}{"membership_id": 1}); isErr {
		t.Error("get_membership failed")
	}
	if _, isErr := invokeTool(t, h["get_membership"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing membership_id")
	}
	if _, isErr := invokeTool(t, h["create_membership"], map[string]interface{}{
		"project_id": "alpha", "user_id": 1, "role_ids": []interface{}{3.0, 4.0},
	}); isErr {
		t.Error("create_membership failed")
	}
	if _, isErr := invokeTool(t, h["create_membership"], map[string]interface{}{"user_id": 1, "role_ids": []interface{}{3.0}}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["create_membership"], map[string]interface{}{"project_id": "alpha", "role_ids": []interface{}{3.0}}); !isErr {
		t.Error("expected error for missing user_id")
	}
	if _, isErr := invokeTool(t, h["create_membership"], map[string]interface{}{"project_id": "alpha", "user_id": 1}); !isErr {
		t.Error("expected error for missing role_ids")
	}
	if _, isErr := invokeTool(t, h["update_membership"], map[string]interface{}{
		"membership_id": 1, "role_ids": []interface{}{3.0},
	}); isErr {
		t.Error("update_membership failed")
	}
	if _, isErr := invokeTool(t, h["update_membership"], map[string]interface{}{"role_ids": []interface{}{3.0}}); !isErr {
		t.Error("expected error for missing membership_id")
	}
	if _, isErr := invokeTool(t, h["update_membership"], map[string]interface{}{"membership_id": 1}); !isErr {
		t.Error("expected error for missing role_ids")
	}
	if _, isErr := invokeTool(t, h["delete_membership"], map[string]interface{}{"membership_id": 1}); isErr {
		t.Error("delete_membership failed")
	}
	if _, isErr := invokeTool(t, h["delete_membership"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing membership_id")
	}
}

// --- Wiki ------------------------------------------------------------------

func TestWikiTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/projects/alpha/wiki/index.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"wiki_pages": []map[string]interface{}{{"title": "Home", "version": 1}},
			})
		case r.URL.Path == "/projects/alpha/wiki/Home.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"wiki_page": map[string]interface{}{"title": "Home", "text": "Hello", "version": 1},
			})
		case r.URL.Path == "/projects/alpha/wiki/Home.json" && (r.Method == http.MethodPut || r.Method == http.MethodPost):
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"wiki_page": map[string]interface{}{"title": "Home", "text": "Hi", "version": 2},
			})
		case r.URL.Path == "/projects/alpha/wiki/Home.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := api.NewClient(srv.URL, "k")
	h := registerAndCapture(t, registerWikiTools, client)

	if _, isErr := invokeTool(t, h["list_wiki_pages"], map[string]interface{}{"project_id": "alpha"}); isErr {
		t.Error("list_wiki_pages failed")
	}
	if _, isErr := invokeTool(t, h["list_wiki_pages"], map[string]interface{}{}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["get_wiki_page"], map[string]interface{}{"project_id": "alpha", "title": "Home"}); isErr {
		t.Error("get_wiki_page failed")
	}
	if _, isErr := invokeTool(t, h["get_wiki_page"], map[string]interface{}{"title": "Home"}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["get_wiki_page"], map[string]interface{}{"project_id": "alpha"}); !isErr {
		t.Error("expected error for missing title")
	}
	if _, isErr := invokeTool(t, h["create_or_update_wiki_page"], map[string]interface{}{
		"project_id": "alpha", "title": "Home", "text": "Hi",
	}); isErr {
		t.Error("create_or_update_wiki_page failed")
	}
	if _, isErr := invokeTool(t, h["create_or_update_wiki_page"], map[string]interface{}{"title": "Home", "text": "x"}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["create_or_update_wiki_page"], map[string]interface{}{"project_id": "alpha", "text": "x"}); !isErr {
		t.Error("expected error for missing title")
	}
	if _, isErr := invokeTool(t, h["create_or_update_wiki_page"], map[string]interface{}{"project_id": "alpha", "title": "Home"}); !isErr {
		t.Error("expected error for missing text")
	}
	if _, isErr := invokeTool(t, h["delete_wiki_page"], map[string]interface{}{"project_id": "alpha", "title": "Home"}); isErr {
		t.Error("delete_wiki_page failed")
	}
	if _, isErr := invokeTool(t, h["delete_wiki_page"], map[string]interface{}{"title": "Home"}); !isErr {
		t.Error("expected error for missing project_id")
	}
	if _, isErr := invokeTool(t, h["delete_wiki_page"], map[string]interface{}{"project_id": "alpha"}); !isErr {
		t.Error("expected error for missing title")
	}
}

// --- getIntSliceArg --------------------------------------------------------

func TestGetIntSliceArg(t *testing.T) {
	args := map[string]interface{}{
		"numbers": []interface{}{1.0, 2.0, 3.0},
		"ints":    []interface{}{int(4), int(5)},
		"mixed":   []interface{}{1.0, "bad", 2.0},
		"empty":   []interface{}{},
		"not":     "string",
	}
	if got := getIntSliceArg(args, "numbers"); len(got) != 3 || got[0] != 1 {
		t.Errorf("numbers: %v", got)
	}
	if got := getIntSliceArg(args, "ints"); len(got) != 2 || got[0] != 4 {
		t.Errorf("ints: %v", got)
	}
	if got := getIntSliceArg(args, "mixed"); len(got) != 2 {
		t.Errorf("mixed: expected 2 entries, got %v", got)
	}
	if got := getIntSliceArg(args, "empty"); len(got) != 0 {
		t.Errorf("empty: %v", got)
	}
	if got := getIntSliceArg(args, "not"); got != nil {
		t.Errorf("not: expected nil, got %v", got)
	}
	if got := getIntSliceArg(args, "missing"); got != nil {
		t.Errorf("missing: expected nil, got %v", got)
	}
}
