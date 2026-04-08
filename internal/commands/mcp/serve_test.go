package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nbifrye/rmn/internal/api"
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

// Ensure unused import doesn't fail build
var _ = mcplib.MethodToolsCall
