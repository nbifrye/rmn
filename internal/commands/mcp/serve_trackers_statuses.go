package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerTrackerTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_trackers",
		Title:       "List Redmine Trackers",
		Description: "List all Redmine trackers (issue types). Returns a JSON array with id/name fields. Tracker IDs are specific to each Redmine instance and are needed for filtering issues or creating issues with a specific tracker.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		trackers, err := client.ListTrackers(ctx)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Trackers []api.IdName `json:"trackers"`
		}{Trackers: trackers})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}

func registerStatusTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_issue_statuses",
		Title:       "List Redmine Issue Statuses",
		Description: "List all Redmine issue statuses. Returns a JSON array with id/name/is_closed fields. Status IDs are specific to each Redmine instance and are needed for updating issue status.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		statuses, err := client.ListStatuses(ctx)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			IssueStatuses []api.IssueStatus `json:"issue_statuses"`
		}{IssueStatuses: statuses})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
