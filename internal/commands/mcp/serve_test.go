package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
)

func setupTestServer(t *testing.T) (*server.MCPServer, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/issues.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"issues": []map[string]interface{}{
					{"id": 1, "subject": "Test issue"},
				},
				"total_count": 1,
			})
		case r.URL.Path == "/issues/42.json" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"issue": map[string]interface{}{
					"id": 42, "subject": "Found issue",
				},
			})
		case r.URL.Path == "/issues.json" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"issue": map[string]interface{}{
					"id": 100, "subject": "Created issue",
				},
			})
		case r.URL.Path == "/issues/42.json" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/issues/42.json" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errors":["Not found"]}`))
		}
	}))

	client := api.NewClient(srv.URL, "test-key")
	mcpServer := server.NewMCPServer("test", "0.1.0")
	registerTools(mcpServer, client)

	// Initialize the server first
	initMsg := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	mcpServer.HandleMessage(context.Background(), json.RawMessage(initMsg))

	return mcpServer, srv
}

func callTool(t *testing.T, s *server.MCPServer, name string, args map[string]interface{}) (string, bool) {
	t.Helper()

	argsJSON, _ := json.Marshal(args)
	msg := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":%q,"arguments":%s}}`, name, string(argsJSON))

	resp := s.HandleMessage(context.Background(), json.RawMessage(msg))

	// Parse the response
	respData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var rpcResp struct {
		Result struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			IsError bool `json:"isError"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respData, &rpcResp); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, string(respData))
	}

	if rpcResp.Error != nil {
		return rpcResp.Error.Message, true
	}

	if len(rpcResp.Result.Content) == 0 {
		t.Fatal("empty result content")
	}

	return rpcResp.Result.Content[0].Text, rpcResp.Result.IsError
}

func TestListIssuesHandler(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "list_issues", map[string]interface{}{
		"project_id": "test",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}

	var resp struct {
		Issues     []api.Issue `json:"issues"`
		TotalCount int         `json:"total_count"`
	}
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, text)
	}
	if resp.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", resp.TotalCount)
	}
}

func TestGetIssueHandler(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "get_issue", map[string]interface{}{
		"issue_id": 42,
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}

	var issue api.Issue
	if err := json.Unmarshal([]byte(text), &issue); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, text)
	}
	if issue.ID != 42 {
		t.Errorf("expected issue ID 42, got %d", issue.ID)
	}
}

func TestGetIssueHandler_MissingID(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "get_issue", map[string]interface{}{})
	if !isErr {
		t.Error("expected error result")
	}
	if text != "issue_id is required" {
		t.Errorf("unexpected error: %s", text)
	}
}

func TestCreateIssueHandler(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "create_issue", map[string]interface{}{
		"project_id": "test",
		"subject":    "New issue",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}

	var issue api.Issue
	if err := json.Unmarshal([]byte(text), &issue); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, text)
	}
	if issue.ID != 100 {
		t.Errorf("expected issue ID 100, got %d", issue.ID)
	}
}

func TestCreateIssueHandler_MissingRequired(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	// Missing subject
	text, isErr := callTool(t, s, "create_issue", map[string]interface{}{
		"project_id": "test",
	})
	if !isErr {
		t.Errorf("expected error for missing subject, got: %s", text)
	}

	// Missing project_id
	text, isErr = callTool(t, s, "create_issue", map[string]interface{}{
		"subject": "test",
	})
	if !isErr {
		t.Errorf("expected error for missing project_id, got: %s", text)
	}
}

func TestUpdateIssueHandler(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "update_issue", map[string]interface{}{
		"issue_id": 42,
		"notes":    "Updated via MCP",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, text)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
}

func TestDeleteIssueHandler(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "delete_issue", map[string]interface{}{
		"issue_id": 42,
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}

	var resp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, text)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
}

func TestDeleteIssueHandler_MissingID(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	_, isErr := callTool(t, s, "delete_issue", map[string]interface{}{})
	if !isErr {
		t.Error("expected error result")
	}
}

func TestGetStringPtrArg_NonStringType(t *testing.T) {
	args := map[string]interface{}{
		"str":     "hello",
		"number":  42.0,
		"nil_val": nil,
		"bool":    true,
	}

	// String type should return pointer to the string
	if v := getStringPtrArg(args, "str"); v == nil || *v != "hello" {
		t.Errorf("expected pointer to 'hello', got %v", v)
	}

	// Non-string type (number) should return nil, not &""
	if v := getStringPtrArg(args, "number"); v != nil {
		t.Errorf("expected nil for number type, got %q", *v)
	}

	// Nil value should return nil
	if v := getStringPtrArg(args, "nil_val"); v != nil {
		t.Errorf("expected nil for nil value, got %q", *v)
	}

	// Absent key should return nil
	if v := getStringPtrArg(args, "missing"); v != nil {
		t.Errorf("expected nil for missing key, got %q", *v)
	}

	// Bool type should return nil
	if v := getStringPtrArg(args, "bool"); v != nil {
		t.Errorf("expected nil for bool type, got %q", *v)
	}
}

func TestToolAnnotations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := api.NewClient(srv.URL, "test-key")
	mcpServer := server.NewMCPServer("test", "0.1.0")
	registerTools(mcpServer, client)

	// Initialize the server
	initMsg := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	mcpServer.HandleMessage(context.Background(), json.RawMessage(initMsg))

	// List tools to verify annotations
	listMsg := `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`
	resp := mcpServer.HandleMessage(context.Background(), json.RawMessage(listMsg))

	respData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var rpcResp struct {
		Result struct {
			Tools []struct {
				Name        string `json:"name"`
				Annotations struct {
					Title           string `json:"title"`
					ReadOnlyHint    *bool  `json:"readOnlyHint"`
					DestructiveHint *bool  `json:"destructiveHint"`
					IdempotentHint  *bool  `json:"idempotentHint"`
					OpenWorldHint   *bool  `json:"openWorldHint"`
				} `json:"annotations"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(respData, &rpcResp); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, string(respData))
	}

	tests := []struct {
		name            string
		wantReadOnly    bool
		wantDestructive bool
		wantIdempotent  bool
	}{
		{"list_issues", true, false, true},
		{"get_issue", true, false, true},
		{"create_issue", false, false, false},
		{"update_issue", false, false, true},
		{"delete_issue", false, true, true},
	}

	toolMap := make(map[string]struct {
		ReadOnlyHint    *bool
		DestructiveHint *bool
		IdempotentHint  *bool
	})
	for _, tool := range rpcResp.Result.Tools {
		toolMap[tool.Name] = struct {
			ReadOnlyHint    *bool
			DestructiveHint *bool
			IdempotentHint  *bool
		}{
			ReadOnlyHint:    tool.Annotations.ReadOnlyHint,
			DestructiveHint: tool.Annotations.DestructiveHint,
			IdempotentHint:  tool.Annotations.IdempotentHint,
		}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tool, ok := toolMap[tc.name]
			if !ok {
				t.Fatalf("tool %q not found in registered tools", tc.name)
			}
			if tool.ReadOnlyHint == nil || *tool.ReadOnlyHint != tc.wantReadOnly {
				t.Errorf("readOnlyHint: got %v, want %v", tool.ReadOnlyHint, tc.wantReadOnly)
			}
			if tool.DestructiveHint == nil || *tool.DestructiveHint != tc.wantDestructive {
				t.Errorf("destructiveHint: got %v, want %v", tool.DestructiveHint, tc.wantDestructive)
			}
			if tool.IdempotentHint == nil || *tool.IdempotentHint != tc.wantIdempotent {
				t.Errorf("idempotentHint: got %v, want %v", tool.IdempotentHint, tc.wantIdempotent)
			}
		})
	}
}

func TestNewCmdMcp(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
		APIClient: func() (*api.Client, error) {
			return nil, nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewCmdMcp(f, "test")
	if cmd.Use != "mcp" {
		t.Errorf("expected Use 'mcp', got %q", cmd.Use)
	}

	// Verify serve subcommand exists
	found := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "serve" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'serve' subcommand not found")
	}
}

func TestNewCmdServe_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := &cmdutil.Factory{
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

	// Replace serveStdioFunc to avoid blocking on stdin
	origServe := serveStdioFunc
	serveStdioFunc = func(s *server.MCPServer) error { return nil }
	defer func() { serveStdioFunc = origServe }()

	cmd := newCmdServe(f, "test")
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCmdServe_APIClientError(t *testing.T) {
	f := &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
		APIClient: func() (*api.Client, error) {
			return nil, fmt.Errorf("not configured")
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := newCmdServe(f, "test")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when API client fails")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("not configured")) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetArgs_NonMap(t *testing.T) {
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = "not a map"
	args := getArgs(req)
	if len(args) != 0 {
		t.Errorf("expected empty map for non-map arguments, got %v", args)
	}
}

func TestGetArgs_NilArguments(t *testing.T) {
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = nil
	args := getArgs(req)
	if len(args) != 0 {
		t.Errorf("expected empty map for nil arguments, got %v", args)
	}
}

func TestGetIntArg_AllTypes(t *testing.T) {
	args := map[string]interface{}{
		"float":          42.0,
		"int":            int(10),
		"string":         "7",
		"invalid_string": "abc",
		"nil_val":        nil,
		"bool":           true,
	}

	if v := getIntArg(args, "float"); v != 42 {
		t.Errorf("float: expected 42, got %d", v)
	}
	if v := getIntArg(args, "int"); v != 10 {
		t.Errorf("int: expected 10, got %d", v)
	}
	if v := getIntArg(args, "string"); v != 7 {
		t.Errorf("string: expected 7, got %d", v)
	}
	if v := getIntArg(args, "invalid_string"); v != 0 {
		t.Errorf("invalid_string: expected 0, got %d", v)
	}
	if v := getIntArg(args, "nil_val"); v != 0 {
		t.Errorf("nil_val: expected 0, got %d", v)
	}
	if v := getIntArg(args, "missing"); v != 0 {
		t.Errorf("missing: expected 0, got %d", v)
	}
	if v := getIntArg(args, "bool"); v != 0 {
		t.Errorf("bool: expected 0, got %d", v)
	}
}

func TestGetIntPtrArg_AllTypes(t *testing.T) {
	args := map[string]interface{}{
		"float":          42.0,
		"int":            int(10),
		"string":         "7",
		"invalid_string": "abc",
		"nil_val":        nil,
		"bool":           true,
	}

	if v := getIntPtrArg(args, "float"); v == nil || *v != 42 {
		t.Errorf("float: expected *42, got %v", v)
	}
	if v := getIntPtrArg(args, "int"); v == nil || *v != 10 {
		t.Errorf("int: expected *10, got %v", v)
	}
	if v := getIntPtrArg(args, "string"); v == nil || *v != 7 {
		t.Errorf("string: expected *7, got %v", v)
	}
	if v := getIntPtrArg(args, "invalid_string"); v != nil {
		t.Errorf("invalid_string: expected nil, got %v", *v)
	}
	if v := getIntPtrArg(args, "nil_val"); v != nil {
		t.Errorf("nil_val: expected nil, got %v", *v)
	}
	if v := getIntPtrArg(args, "missing"); v != nil {
		t.Errorf("missing: expected nil, got %v", *v)
	}
	if v := getIntPtrArg(args, "bool"); v != nil {
		t.Errorf("bool: expected nil, got %v", *v)
	}
}

func TestToJSON_Error(t *testing.T) {
	// Channels cannot be marshaled to JSON
	_, err := toJSON(make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable value")
	}
}

// setupErrorServer creates an MCP server backed by an HTTP server that always
// returns the given status code and response body, for testing API error handling.
func setupErrorServer(t *testing.T, statusCode int, body string) (*server.MCPServer, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))

	client := api.NewClient(srv.URL, "test-key")
	mcpServer := server.NewMCPServer("test", "0.1.0")
	registerTools(mcpServer, client)

	initMsg := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	mcpServer.HandleMessage(context.Background(), json.RawMessage(initMsg))

	return mcpServer, srv
}

func TestListIssuesHandler_APIError(t *testing.T) {
	s, srv := setupErrorServer(t, http.StatusInternalServerError, `{"errors":["Server error"]}`)
	defer srv.Close()

	text, isErr := callTool(t, s, "list_issues", map[string]interface{}{})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
}

func TestGetIssueHandler_APIError(t *testing.T) {
	s, srv := setupErrorServer(t, http.StatusNotFound, `{"errors":["Not found"]}`)
	defer srv.Close()

	text, isErr := callTool(t, s, "get_issue", map[string]interface{}{"issue_id": 999})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
}

func TestCreateIssueHandler_APIError(t *testing.T) {
	s, srv := setupErrorServer(t, http.StatusUnprocessableEntity, `{"errors":["Subject cannot be blank"]}`)
	defer srv.Close()

	text, isErr := callTool(t, s, "create_issue", map[string]interface{}{
		"project_id": "test",
		"subject":    "test",
	})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
}

func TestCreateIssueHandler_NumericProjectID(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "create_issue", map[string]interface{}{
		"project_id": "1",
		"subject":    "Numeric project",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
}

func TestUpdateIssueHandler_MissingID(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "update_issue", map[string]interface{}{})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
	if text != "issue_id is required" {
		t.Errorf("unexpected error: %s", text)
	}
}

func TestUpdateIssueHandler_APIError(t *testing.T) {
	s, srv := setupErrorServer(t, http.StatusNotFound, `{"errors":["Not found"]}`)
	defer srv.Close()

	text, isErr := callTool(t, s, "update_issue", map[string]interface{}{
		"issue_id": 42,
		"notes":    "test",
	})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
}

func TestDeleteIssueHandler_APIError(t *testing.T) {
	s, srv := setupErrorServer(t, http.StatusNotFound, `{"errors":["Not found"]}`)
	defer srv.Close()

	text, isErr := callTool(t, s, "delete_issue", map[string]interface{}{
		"issue_id": 999,
	})
	if !isErr {
		t.Errorf("expected error result, got: %s", text)
	}
}

func TestUpdateIssueHandler_WithAllFields(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "update_issue", map[string]interface{}{
		"issue_id":       42,
		"subject":        "Updated subject",
		"description":    "New description",
		"status_id":      3,
		"priority_id":    2,
		"assigned_to_id": 5,
		"notes":          "A note",
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
}

func TestCreateIssueHandler_WithAllFields(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "create_issue", map[string]interface{}{
		"project_id":     "test",
		"subject":        "Full issue",
		"description":    "Desc",
		"tracker_id":     1,
		"priority_id":    2,
		"assigned_to_id": 3,
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
}

func TestListIssuesHandler_WithAllParams(t *testing.T) {
	s, srv := setupTestServer(t)
	defer srv.Close()

	text, isErr := callTool(t, s, "list_issues", map[string]interface{}{
		"project_id":      "test",
		"status_id":       "open",
		"assigned_to_id":  "me",
		"tracker_id":      1,
		"limit":           10,
		"offset":          5,
	})
	if isErr {
		t.Fatalf("unexpected error: %s", text)
	}
}

// Ensure unused import doesn't fail build
var _ = mcplib.MethodToolsCall
