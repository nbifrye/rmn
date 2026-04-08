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

// Ensure unused import doesn't fail build
var _ = mcplib.MethodToolsCall
