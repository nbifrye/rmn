package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerVersionTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_versions",
		Title:       "List Redmine Versions",
		Description: "List all versions (milestones) for a given Redmine project. Returns a JSON object with a 'versions' array and 'total_count'.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID (numeric) or identifier (string)"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		versions, total, err := client.ListVersions(ctx, projectID)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Versions   []api.Version `json:"versions"`
			TotalCount int           `json:"total_count"`
		}{Versions: versions, TotalCount: total})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_version",
		Title:       "Get Redmine Version",
		Description: "Get full details of a specific Redmine version by numeric ID.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"version_id": map[string]any{"type": "number", "description": "The numeric version ID"},
			},
			"required": []string{"version_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "version_id")
		if id <= 0 {
			return errResult("version_id must be a positive integer"), nil
		}
		version, err := client.GetVersion(ctx, id)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(version)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "create_version",
		Title:       "Create Redmine Version",
		Description: "Create a new version (milestone) in a Redmine project. Returns the created version with its assigned ID.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id":      map[string]any{"type": "string", "description": "Project ID or identifier"},
				"name":            map[string]any{"type": "string", "description": "Version name"},
				"status":          map[string]any{"type": "string", "description": "Status: 'open', 'locked', or 'closed'"},
				"sharing":         map[string]any{"type": "string", "description": "Sharing: 'none', 'descendants', 'hierarchy', 'tree', 'system'"},
				"due_date":        map[string]any{"type": "string", "description": "Due date in YYYY-MM-DD format"},
				"description":     map[string]any{"type": "string", "description": "Version description"},
				"wiki_page_title": map[string]any{"type": "string", "description": "Associated wiki page title"},
			},
			"required": []string{"project_id", "name"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		params := api.VersionCreateParams{
			Name:          getStringArg(args, "name"),
			Status:        getStringArg(args, "status"),
			Sharing:       getStringArg(args, "sharing"),
			DueDate:       getStringArg(args, "due_date"),
			Description:   getStringArg(args, "description"),
			WikiPageTitle: getStringArg(args, "wiki_page_title"),
		}
		if params.Name == "" {
			return errResult("name is required"), nil
		}
		version, err := client.CreateVersion(ctx, projectID, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(version)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "update_version",
		Title:       "Update Redmine Version",
		Description: "Update an existing Redmine version. Only provided fields are changed.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"version_id":      map[string]any{"type": "number", "description": "The numeric version ID"},
				"name":            map[string]any{"type": "string"},
				"status":          map[string]any{"type": "string", "description": "'open', 'locked', or 'closed'"},
				"sharing":         map[string]any{"type": "string"},
				"due_date":        map[string]any{"type": "string", "description": "YYYY-MM-DD"},
				"description":     map[string]any{"type": "string"},
				"wiki_page_title": map[string]any{"type": "string"},
			},
			"required": []string{"version_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "version_id")
		if id <= 0 {
			return errResult("version_id must be a positive integer"), nil
		}
		params := api.VersionUpdateParams{
			Name:          getStringPtrArg(args, "name"),
			Status:        getStringPtrArg(args, "status"),
			Sharing:       getStringPtrArg(args, "sharing"),
			DueDate:       getStringPtrArg(args, "due_date"),
			Description:   getStringPtrArg(args, "description"),
			WikiPageTitle: getStringPtrArg(args, "wiki_page_title"),
		}
		if err := client.UpdateVersion(ctx, id, params); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated version #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "delete_version",
		Title:       "Delete Redmine Version",
		Description: "Permanently delete a Redmine version. This action cannot be undone. No confirmation is requested.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"version_id": map[string]any{"type": "number", "description": "The numeric version ID"},
			},
			"required": []string{"version_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "version_id")
		if id <= 0 {
			return errResult("version_id must be a positive integer"), nil
		}
		if err := client.DeleteVersion(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted version #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
