package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerWikiTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_wiki_pages",
		Title:       "List Redmine Wiki Pages",
		Description: "List all wiki pages (index) for a Redmine project. Returns a JSON object with a 'wiki_pages' array. Page content is not included; use get_wiki_page to fetch content.",
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
			},
			"required": []string{"project_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		pages, err := client.ListWikiPages(ctx, projectID)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			WikiPages []api.WikiPage `json:"wiki_pages"`
		}{WikiPages: pages})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_wiki_page",
		Title:       "Get Redmine Wiki Page",
		Description: "Get the full content of a specific Redmine wiki page by title. Optionally fetch a specific historical version.",
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
				"title":      map[string]any{"type": "string", "description": "Wiki page title"},
				"version":    map[string]any{"type": "number", "description": "Optional historical version number"},
			},
			"required": []string{"project_id", "title"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		title := getStringArg(args, "title")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		if title == "" {
			return errResult("title is required"), nil
		}
		version := getIntArg(args, "version")
		page, err := client.GetWikiPage(ctx, projectID, title, version)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(page)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "create_or_update_wiki_page",
		Title:       "Create or Update Redmine Wiki Page",
		Description: "Create a new Redmine wiki page or update an existing one. Redmine uses PUT for both operations. Provide the full page text.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
				"title":      map[string]any{"type": "string", "description": "Wiki page title"},
				"text":       map[string]any{"type": "string", "description": "Full page text (Textile or Markdown depending on instance config)"},
				"comments":   map[string]any{"type": "string", "description": "Edit comments/summary"},
			},
			"required": []string{"project_id", "title", "text"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		title := getStringArg(args, "title")
		text := getStringArg(args, "text")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		if title == "" {
			return errResult("title is required"), nil
		}
		if text == "" {
			return errResult("text is required"), nil
		}
		params := api.WikiPageCreateParams{
			Text:     text,
			Comments: getStringArg(args, "comments"),
		}
		page, err := client.CreateWikiPage(ctx, projectID, title, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		result, err := toJSON(page)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(result), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "delete_wiki_page",
		Title:       "Delete Redmine Wiki Page",
		Description: "Permanently delete a Redmine wiki page. This also deletes all historical versions and cannot be undone. No confirmation is requested.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{"type": "string", "description": "Project ID or identifier"},
				"title":      map[string]any{"type": "string", "description": "Wiki page title"},
			},
			"required": []string{"project_id", "title"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		projectID := getStringArg(args, "project_id")
		title := getStringArg(args, "title")
		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		if title == "" {
			return errResult("title is required"), nil
		}
		if err := client.DeleteWikiPage(ctx, projectID, title); err != nil {
			return errResult(err.Error()), nil
		}
		result, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted wiki page %q in project %s", title, projectID)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(result), nil
	})
}
