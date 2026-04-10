package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
)

func registerTimeEntryTools(s *mcp.Server, client *api.Client) {
	addTool(s, &mcp.Tool{
		Name:        "list_time_entries",
		Title:       "List Redmine Time Entries",
		Description: "List Redmine time entries with optional filters. Returns a JSON object with a 'time_entries' array and 'total_count'. Use 'from'/'to' for date ranges (YYYY-MM-DD).",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id":  map[string]any{"type": "string", "description": "Filter by project ID or identifier"},
				"issue_id":    map[string]any{"type": "number", "description": "Filter by issue ID"},
				"user_id":     map[string]any{"type": "number", "description": "Filter by user ID"},
				"spent_on":    map[string]any{"type": "string", "description": "Exact spent date YYYY-MM-DD"},
				"from":        map[string]any{"type": "string", "description": "Range start YYYY-MM-DD"},
				"to":          map[string]any{"type": "string", "description": "Range end YYYY-MM-DD"},
				"activity_id": map[string]any{"type": "number", "description": "Filter by activity ID"},
				"limit":       map[string]any{"type": "number", "description": "Max number of results"},
				"offset":      map[string]any{"type": "number", "description": "Pagination offset"},
			},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		params := api.TimeEntryListParams{
			ProjectID:  getStringArg(args, "project_id"),
			IssueID:    getIntArg(args, "issue_id"),
			UserID:     getIntArg(args, "user_id"),
			SpentOn:    getStringArg(args, "spent_on"),
			From:       getStringArg(args, "from"),
			To:         getStringArg(args, "to"),
			ActivityID: getIntArg(args, "activity_id"),
			Limit:      getIntArg(args, "limit"),
			Offset:     getIntArg(args, "offset"),
		}
		entries, total, err := client.ListTimeEntries(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			TimeEntries []api.TimeEntry `json:"time_entries"`
			TotalCount  int             `json:"total_count"`
		}{TimeEntries: entries, TotalCount: total})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "get_time_entry",
		Title:       "Get Redmine Time Entry",
		Description: "Get full details of a specific Redmine time entry by numeric ID.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"time_entry_id": map[string]any{"type": "number", "description": "The numeric time entry ID"},
			},
			"required": []string{"time_entry_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "time_entry_id")
		if id <= 0 {
			return errResult("time_entry_id must be a positive integer"), nil
		}
		entry, err := client.GetTimeEntry(ctx, id)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(entry)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "create_time_entry",
		Title:       "Create Redmine Time Entry",
		Description: "Log time spent on an issue or project. One of 'issue_id' or 'project_id' is required, along with 'hours'. Returns the created time entry.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"issue_id":    map[string]any{"type": "number", "description": "Issue ID (either this or project_id is required)"},
				"project_id":  map[string]any{"type": "string", "description": "Project ID or identifier (either this or issue_id is required)"},
				"hours":       map[string]any{"type": "number", "description": "Hours spent (required)"},
				"activity_id": map[string]any{"type": "number", "description": "Activity ID (values are specific to your Redmine instance)"},
				"spent_on":    map[string]any{"type": "string", "description": "Date YYYY-MM-DD (defaults to today)"},
				"comments":    map[string]any{"type": "string", "description": "Optional comments"},
			},
			"required": []string{"hours"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		params := api.TimeEntryCreateParams{
			IssueID:    getIntArg(args, "issue_id"),
			ProjectID:  getStringArg(args, "project_id"),
			SpentOn:    getStringArg(args, "spent_on"),
			Hours:      getFloat64Arg(args, "hours"),
			ActivityID: getIntArg(args, "activity_id"),
			Comments:   getStringArg(args, "comments"),
		}
		if params.Hours <= 0 {
			return errResult("hours must be a positive number"), nil
		}
		if params.IssueID == 0 && params.ProjectID == "" {
			return errResult("either issue_id or project_id is required"), nil
		}
		entry, err := client.CreateTimeEntry(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(entry)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "update_time_entry",
		Title:       "Update Redmine Time Entry",
		Description: "Update an existing Redmine time entry. Only provided fields are changed.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"time_entry_id": map[string]any{"type": "number", "description": "The numeric time entry ID"},
				"issue_id":      map[string]any{"type": "number"},
				"project_id":    map[string]any{"type": "string"},
				"hours":         map[string]any{"type": "number"},
				"activity_id":   map[string]any{"type": "number"},
				"spent_on":      map[string]any{"type": "string", "description": "YYYY-MM-DD"},
				"comments":      map[string]any{"type": "string"},
			},
			"required": []string{"time_entry_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "time_entry_id")
		if id <= 0 {
			return errResult("time_entry_id must be a positive integer"), nil
		}
		params := api.TimeEntryUpdateParams{
			IssueID:    getIntPtrArg(args, "issue_id"),
			ProjectID:  getStringPtrArg(args, "project_id"),
			SpentOn:    getStringPtrArg(args, "spent_on"),
			Hours:      getFloat64PtrArg(args, "hours"),
			ActivityID: getIntPtrArg(args, "activity_id"),
			Comments:   getStringPtrArg(args, "comments"),
		}
		if err := client.UpdateTimeEntry(ctx, id, params); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated time entry #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})

	addTool(s, &mcp.Tool{
		Name:        "delete_time_entry",
		Title:       "Delete Redmine Time Entry",
		Description: "Permanently delete a Redmine time entry. This action cannot be undone. No confirmation is requested.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"time_entry_id": map[string]any{"type": "number", "description": "The numeric time entry ID"},
			},
			"required": []string{"time_entry_id"},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)
		id := getIntArg(args, "time_entry_id")
		if id <= 0 {
			return errResult("time_entry_id must be a positive integer"), nil
		}
		if err := client.DeleteTimeEntry(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}
		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted time entry #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	})
}
