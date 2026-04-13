package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerUserTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_users",
		Title:       "List Redmine Users",
		Description: "List Redmine users with optional filters. Returns a JSON object with a 'users' array and 'total_count'. Requires admin privileges on the Redmine instance.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"status":   map[string]any{"type": "number", "description": "Filter by status: 1=active, 2=registered, 3=locked"},
				"name":     map[string]any{"type": "string", "description": "Filter by name/login substring"},
				"group_id": map[string]any{"type": "number", "description": "Filter by group ID"},
				"limit":    map[string]any{"type": "number", "description": "Max number of results"},
				"offset":   map[string]any{"type": "number", "description": "Pagination offset"},
			},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		params := api.UserListParams{
			Status:  getIntArg(args, "status"),
			Name:    getStringArg(args, "name"),
			GroupID: getIntArg(args, "group_id"),
			Limit:   getIntArg(args, "limit"),
			Offset:  getIntArg(args, "offset"),
		}
		users, total, err := client.ListUsers(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Users      []api.User `json:"users"`
			TotalCount int        `json:"total_count"`
		}{Users: users, TotalCount: total})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_user",
		Title:       "Get Redmine User",
		Description: "Get full details of a specific Redmine user by numeric ID.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"user_id": map[string]any{"type": "number", "description": "The numeric user ID"},
			},
			"required": []string{"user_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "user_id")
		if id <= 0 {
			return errResult("user_id must be a positive integer"), nil
		}
		user, err := client.GetUser(ctx, id)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(user)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_current_user",
		Title:       "Get Current Redmine User",
		Description: "Get details of the user whose API key is used for the current session.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user, err := client.GetCurrentUser(ctx)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(user)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
