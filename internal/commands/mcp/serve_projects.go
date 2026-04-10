package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerProjectTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_projects",
		Title:       "List Redmine Projects",
		Description: "List Redmine projects with optional status filter. Returns a JSON object with a 'projects' array and 'total_count'. Use 'offset' and 'limit' for pagination.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"status": map[string]any{"type": "string", "description": "Filter by status: 'active', 'closed', or 'archived'"},
				"limit":  map[string]any{"type": "number", "description": "Max number of results to return"},
				"offset": map[string]any{"type": "number", "description": "Pagination offset"},
			},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		params := api.ProjectListParams{
			Status: getStringArg(args, "status"),
			Limit:  getIntArg(args, "limit"),
			Offset: getIntArg(args, "offset"),
		}
		projects, total, err := client.ListProjects(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Projects   []api.Project `json:"projects"`
			TotalCount int           `json:"total_count"`
		}{Projects: projects, TotalCount: total})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_project",
		Title:       "Get Redmine Project",
		Description: "Get full details of a Redmine project by ID or identifier. Use 'include' to fetch associations like trackers, issue_categories, time_entry_activities.",
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
				"include":    map[string]any{"type": "string", "description": "Comma-separated list of associations to include"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getStringArg(args, "project_id")
		if id == "" {
			return errResult("project_id is required"), nil
		}
		var include []string
		if inc := getStringArg(args, "include"); inc != "" {
			include = strings.Split(inc, ",")
		}
		project, err := client.GetProject(ctx, id, include)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(project)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "create_project",
		Title:       "Create Redmine Project",
		Description: "Create a new Redmine project. Returns the created project with its assigned ID.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "Project name"},
				"identifier":  map[string]any{"type": "string", "description": "Project identifier (URL-safe slug)"},
				"description": map[string]any{"type": "string", "description": "Project description"},
				"homepage":    map[string]any{"type": "string", "description": "Project homepage URL"},
				"is_public":   map[string]any{"type": "boolean", "description": "Whether the project is public"},
				"parent_id":   map[string]any{"type": "number", "description": "Parent project ID"},
			},
			"required": []string{"name", "identifier"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		params := api.ProjectCreateParams{
			Name:        getStringArg(args, "name"),
			Identifier:  getStringArg(args, "identifier"),
			Description: getStringArg(args, "description"),
			Homepage:    getStringArg(args, "homepage"),
			ParentID:    getIntArg(args, "parent_id"),
		}
		if v, ok := args["is_public"].(bool); ok {
			params.IsPublic = v
		}
		if params.Name == "" {
			return errResult("name is required"), nil
		}
		if params.Identifier == "" {
			return errResult("identifier is required"), nil
		}
		project, err := client.CreateProject(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(project)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "update_project",
		Title:       "Update Redmine Project",
		Description: "Update an existing Redmine project. Only provided fields are changed.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id":  map[string]any{"type": "string", "description": "Project ID or identifier"},
				"name":        map[string]any{"type": "string"},
				"description": map[string]any{"type": "string"},
				"homepage":    map[string]any{"type": "string"},
				"is_public":   map[string]any{"type": "boolean"},
				"parent_id":   map[string]any{"type": "number"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getStringArg(args, "project_id")
		if id == "" {
			return errResult("project_id is required"), nil
		}
		params := api.ProjectUpdateParams{
			Name:        getStringPtrArg(args, "name"),
			Description: getStringPtrArg(args, "description"),
			Homepage:    getStringPtrArg(args, "homepage"),
			ParentID:    getIntPtrArg(args, "parent_id"),
		}
		if v, ok := args["is_public"].(bool); ok {
			params.IsPublic = &v
		}
		if err := client.UpdateProject(ctx, id, params); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated project %s", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "archive_project",
		Title:       "Archive Redmine Project",
		Description: "Archive a Redmine project. Archived projects are hidden but not deleted. This is reversible via unarchive_project.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getStringArg(args, "project_id")
		if id == "" {
			return errResult("project_id is required"), nil
		}
		if err := client.ArchiveProject(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Archived project %s", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "unarchive_project",
		Title:       "Unarchive Redmine Project",
		Description: "Unarchive a previously archived Redmine project.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getStringArg(args, "project_id")
		if id == "" {
			return errResult("project_id is required"), nil
		}
		if err := client.UnarchiveProject(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Unarchived project %s", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "delete_project",
		Title:       "Delete Redmine Project",
		Description: "Permanently delete a Redmine project. This action cannot be undone and will delete all issues, time entries, wiki pages, and memberships in the project. No confirmation is requested.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getStringArg(args, "project_id")
		if id == "" {
			return errResult("project_id is required"), nil
		}
		if err := client.DeleteProject(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted project %s", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
