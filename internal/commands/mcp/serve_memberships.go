package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

// getIntSliceArg extracts an []int from the arguments. Accepts JSON arrays of
// numbers or strings (numeric). Returns nil if the key is absent.
func getIntSliceArg(args map[string]interface{}, key string) []int {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]int, 0, len(raw))
	for _, item := range raw {
		switch n := item.(type) {
		case float64:
			out = append(out, int(n))
		case int:
			out = append(out, n)
		}
	}
	return out
}

func registerMembershipTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_memberships",
		Title:       "List Redmine Project Memberships",
		Description: "List memberships (users and groups with roles) for a Redmine project. Returns a JSON object with a 'memberships' array and 'total_count'.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
				"limit":      map[string]any{"type": "number", "description": "Max number of results"},
				"offset":     map[string]any{"type": "number", "description": "Pagination offset"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		params := api.MembershipListParams{
			Limit:  getIntArg(args, "limit"),
			Offset: getIntArg(args, "offset"),
		}
		memberships, total, err := client.ListMemberships(ctx, projectID, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Memberships []api.Membership `json:"memberships"`
			TotalCount  int              `json:"total_count"`
		}{Memberships: memberships, TotalCount: total})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_membership",
		Title:       "Get Redmine Membership",
		Description: "Get full details of a specific Redmine project membership by numeric ID.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"membership_id": map[string]any{"type": "number", "description": "The numeric membership ID"},
			},
			"required": []string{"membership_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "membership_id")
		if id <= 0 {
			return errResult("membership_id must be a positive integer"), nil
		}
		m, err := client.GetMembership(ctx, id)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(m)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "create_membership",
		Title:       "Create Redmine Project Membership",
		Description: "Add a user to a Redmine project with the given roles. Returns the created membership.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
				"user_id":    map[string]any{"type": "number", "description": "User ID to add"},
				"role_ids":   map[string]any{"type": "array", "items": map[string]any{"type": "number"}, "description": "Role IDs to assign to the user"},
			},
			"required": []string{"project_id", "user_id", "role_ids"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		userID := getIntArg(args, "user_id")
		if userID <= 0 {
			return errResult("user_id must be a positive integer"), nil
		}
		roleIDs := getIntSliceArg(args, "role_ids")
		if len(roleIDs) == 0 {
			return errResult("role_ids is required"), nil
		}
		params := api.MembershipCreateParams{
			UserID:  userID,
			RoleIDs: roleIDs,
		}
		m, err := client.CreateMembership(ctx, projectID, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(m)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "update_membership",
		Title:       "Update Redmine Membership",
		Description: "Update the roles of a Redmine project membership. Replaces the existing role assignments.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"membership_id": map[string]any{"type": "number", "description": "The numeric membership ID"},
				"role_ids":      map[string]any{"type": "array", "items": map[string]any{"type": "number"}, "description": "New role IDs"},
			},
			"required": []string{"membership_id", "role_ids"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "membership_id")
		if id <= 0 {
			return errResult("membership_id must be a positive integer"), nil
		}
		roleIDs := getIntSliceArg(args, "role_ids")
		if len(roleIDs) == 0 {
			return errResult("role_ids is required"), nil
		}
		if err := client.UpdateMembership(ctx, id, api.MembershipUpdateParams{RoleIDs: roleIDs}); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated membership #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "delete_membership",
		Title:       "Delete Redmine Membership",
		Description: "Remove a user or group from a Redmine project by deleting the membership. This action cannot be undone. No confirmation is requested.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"membership_id": map[string]any{"type": "number", "description": "The numeric membership ID"},
			},
			"required": []string{"membership_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "membership_id")
		if id <= 0 {
			return errResult("membership_id must be a positive integer"), nil
		}
		if err := client.DeleteMembership(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted membership #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
