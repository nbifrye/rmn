package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdMcp(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server commands",
		Long:  "Model Context Protocol server for exposing Redmine operations to AI agents.",
	}

	cmd.AddCommand(newCmdServe(f, version))
	return cmd
}

// runServerFunc is the function used to start the MCP server. It can be replaced in tests.
var runServerFunc = func(ctx context.Context, s *mcp.Server) error {
	return s.Run(ctx, &mcp.StdioTransport{})
}

func newCmdServe(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server",
		Long:  "Start a stdio-based MCP server that exposes Redmine operations as tools.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			s := mcp.NewServer(&mcp.Implementation{
				Name:    "rmn-redmine",
				Version: version,
			}, nil)
			registerTools(s, client)

			return runServerFunc(cmd.Context(), s)
		},
	}

	return cmd
}

func boolPtr(b bool) *bool { return &b }

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

func toolDefs() []*mcp.Tool {
	return []*mcp.Tool{
		{
			Name:        "list_issues",
			Title:       "List Redmine Issues",
			Description: "List Redmine issues with optional filters. Returns a JSON object with an 'issues' array and 'total_count'. Without filters, returns open issues. Use 'offset' and 'limit' for pagination through large result sets.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(true),
			},
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_id":      map[string]any{"type": "string", "description": "Filter by project ID (numeric) or identifier (string, e.g. 'my-project')"},
					"status_id":       map[string]any{"type": "string", "description": "Filter by status: 'open' (default), 'closed', '*' (all), or a numeric status ID"},
					"assigned_to_id":  map[string]any{"type": "string", "description": "Filter by assignee: 'me' for current user, or a numeric user ID"},
					"tracker_id":      map[string]any{"type": "number", "description": "Filter by tracker ID (values are specific to your Redmine instance)"},
					"sort":            map[string]any{"type": "string", "description": "Sort by column, e.g. 'updated_on:desc', 'priority:asc'"},
					"limit":           map[string]any{"type": "number", "description": "Max number of results to return (default 25, max 100)"},
					"offset":          map[string]any{"type": "number", "description": "Pagination offset for retrieving subsequent pages"},
				},
			},
		},
		{
			Name:        "get_issue",
			Title:       "Get Redmine Issue",
			Description: "Get full details of a specific Redmine issue by ID, including project, tracker, status, priority, author, assignee, description, and timestamps. Use 'include' to fetch associated data like journals (comments) and attachments.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(true),
			},
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issue_id": map[string]any{"type": "number", "description": "The numeric issue ID"},
					"include":  map[string]any{"type": "string", "description": "Comma-separated list of associations to include: journals, attachments, relations, changesets, watchers, children"},
				},
				"required": []string{"issue_id"},
			},
		},
		{
			Name:        "create_issue",
			Title:       "Create Redmine Issue",
			Description: "Create a new Redmine issue. Returns the created issue with its assigned ID.",
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: boolPtr(false),
				OpenWorldHint:   boolPtr(true),
			},
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_id":       map[string]any{"type": "string", "description": "Project ID (numeric) or identifier (string, e.g. 'my-project')"},
					"subject":          map[string]any{"type": "string", "description": "Issue subject/title"},
					"description":      map[string]any{"type": "string", "description": "Detailed issue description (supports Textile markup)"},
					"tracker_id":       map[string]any{"type": "number", "description": "Tracker ID (values are specific to your Redmine instance)"},
					"priority_id":      map[string]any{"type": "number", "description": "Priority ID (values are specific to your Redmine instance)"},
					"assigned_to_id":   map[string]any{"type": "number", "description": "User ID to assign the issue to"},
					"category_id":      map[string]any{"type": "number", "description": "Category ID"},
					"fixed_version_id": map[string]any{"type": "number", "description": "Target version ID"},
					"parent_issue_id":  map[string]any{"type": "number", "description": "Parent issue ID"},
					"start_date":       map[string]any{"type": "string", "description": "Start date in YYYY-MM-DD format"},
					"due_date":         map[string]any{"type": "string", "description": "Due date in YYYY-MM-DD format"},
					"estimated_hours":  map[string]any{"type": "number", "description": "Estimated hours for the issue"},
					"done_ratio":       map[string]any{"type": "number", "description": "Done ratio (0-100)"},
				},
				"required": []string{"project_id", "subject"},
			},
		},
		{
			Name:        "update_issue",
			Title:       "Update Redmine Issue",
			Description: "Update an existing Redmine issue. Only provided fields are changed; omitted fields are left unchanged.",
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(true),
			},
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issue_id":         map[string]any{"type": "number", "description": "The numeric issue ID to update"},
					"subject":          map[string]any{"type": "string", "description": "New subject/title"},
					"description":      map[string]any{"type": "string", "description": "New description (supports Textile markup)"},
					"status_id":        map[string]any{"type": "number", "description": "New status ID (values are specific to your Redmine instance)"},
					"priority_id":      map[string]any{"type": "number", "description": "New priority ID (values are specific to your Redmine instance)"},
					"assigned_to_id":   map[string]any{"type": "number", "description": "New assignee user ID (set to 0 to unassign)"},
					"category_id":      map[string]any{"type": "number", "description": "New category ID"},
					"fixed_version_id": map[string]any{"type": "number", "description": "New target version ID"},
					"parent_issue_id":  map[string]any{"type": "number", "description": "New parent issue ID"},
					"start_date":       map[string]any{"type": "string", "description": "New start date in YYYY-MM-DD format"},
					"due_date":         map[string]any{"type": "string", "description": "New due date in YYYY-MM-DD format"},
					"estimated_hours":  map[string]any{"type": "number", "description": "New estimated hours"},
					"done_ratio":       map[string]any{"type": "number", "description": "New done ratio (0-100)"},
					"notes":            map[string]any{"type": "string", "description": "Add a comment/note to the issue"},
				},
				"required": []string{"issue_id"},
			},
		},
		{
			Name:        "delete_issue",
			Title:       "Delete Redmine Issue",
			Description: "Permanently delete a Redmine issue. This action cannot be undone. No confirmation is requested.",
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: boolPtr(true),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(true),
			},
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issue_id": map[string]any{"type": "number", "description": "The numeric issue ID to delete"},
				},
				"required": []string{"issue_id"},
			},
		},
	}
}

func registerTools(s *mcp.Server, client *api.Client) {
	tools := toolDefs()
	handlers := []mcp.ToolHandler{
		makeListIssuesHandler(client),
		makeGetIssueHandler(client),
		makeCreateIssueHandler(client),
		makeUpdateIssueHandler(client),
		makeDeleteIssueHandler(client),
	}
	for i, tool := range tools {
		s.AddTool(tool, handlers[i])
	}
}

// toJSON marshals a value to indented JSON. It is a variable so tests can replace it.
var toJSON = func(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(data), nil
}

func getArgs(req *mcp.CallToolRequest) map[string]interface{} {
	var m map[string]interface{}
	if err := json.Unmarshal(req.Params.Arguments, &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

func getStringArg(args map[string]interface{}, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getIntArg(args map[string]interface{}, key string) int {
	v, ok := args[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

// getIntPtrArg returns nil if the key is absent, or a pointer to the int value.
func getIntPtrArg(args map[string]interface{}, key string) *int {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	switch n := v.(type) {
	case float64:
		i := int(n)
		return &i
	case int:
		return &n
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return nil
		}
		return &i
	}
	return nil
}

// getStringPtrArg returns nil if the key is absent or not a string, or a pointer to the string value.
func getStringPtrArg(args map[string]interface{}, key string) *string {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return &s
}

// getFloat64Arg returns 0 if the key is absent or not a number.
func getFloat64Arg(args map[string]interface{}, key string) float64 {
	v, ok := args[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	}
	return 0
}

// getFloat64PtrArg returns nil if the key is absent, or a pointer to the float64 value.
func getFloat64PtrArg(args map[string]interface{}, key string) *float64 {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	switch n := v.(type) {
	case float64:
		return &n
	case int:
		f := float64(n)
		return &f
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil
		}
		return &f
	}
	return nil
}

func makeListIssuesHandler(client *api.Client) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		params := api.IssueListParams{
			ProjectID:    getStringArg(args, "project_id"),
			StatusID:     getStringArg(args, "status_id"),
			AssignedToID: getStringArg(args, "assigned_to_id"),
			TrackerID:    getIntArg(args, "tracker_id"),
			Sort:         getStringArg(args, "sort"),
			Limit:        getIntArg(args, "limit"),
			Offset:       getIntArg(args, "offset"),
		}

		issues, total, err := client.ListIssues(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}

		result := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: issues, TotalCount: total}

		text, err := toJSON(result)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	}
}

func makeGetIssueHandler(client *api.Client) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return errResult("issue_id must be a positive integer"), nil
		}

		var include []string
		if inc := getStringArg(args, "include"); inc != "" {
			include = strings.Split(inc, ",")
		}

		issue, err := client.GetIssue(ctx, id, include)
		if err != nil {
			return errResult(err.Error()), nil
		}

		text, err := toJSON(issue)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	}
}

func makeCreateIssueHandler(client *api.Client) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		projectID := getStringArg(args, "project_id")
		subject := getStringArg(args, "subject")

		if projectID == "" {
			return errResult("project_id is required"), nil
		}
		if subject == "" {
			return errResult("subject is required"), nil
		}

		var parsedProjectID interface{} = projectID
		if id, err := strconv.Atoi(projectID); err == nil {
			parsedProjectID = id
		}

		params := api.IssueCreateParams{
			ProjectID:      parsedProjectID,
			Subject:        subject,
			Description:    getStringArg(args, "description"),
			TrackerID:      getIntArg(args, "tracker_id"),
			PriorityID:     getIntArg(args, "priority_id"),
			AssignedToID:   getIntArg(args, "assigned_to_id"),
			CategoryID:     getIntArg(args, "category_id"),
			FixedVersionID: getIntArg(args, "fixed_version_id"),
			ParentIssueID:  getIntArg(args, "parent_issue_id"),
			StartDate:      getStringArg(args, "start_date"),
			DueDate:        getStringArg(args, "due_date"),
			EstimatedHours: getFloat64Arg(args, "estimated_hours"),
			DoneRatio:      getIntArg(args, "done_ratio"),
		}

		issue, err := client.CreateIssue(ctx, params)
		if err != nil {
			return errResult(err.Error()), nil
		}

		text, err := toJSON(issue)
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	}
}

func makeUpdateIssueHandler(client *api.Client) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return errResult("issue_id must be a positive integer"), nil
		}

		params := api.IssueUpdateParams{
			Subject:        getStringPtrArg(args, "subject"),
			Description:    getStringPtrArg(args, "description"),
			StatusID:       getIntPtrArg(args, "status_id"),
			PriorityID:     getIntPtrArg(args, "priority_id"),
			AssignedToID:   getIntPtrArg(args, "assigned_to_id"),
			CategoryID:     getIntPtrArg(args, "category_id"),
			FixedVersionID: getIntPtrArg(args, "fixed_version_id"),
			ParentIssueID:  getIntPtrArg(args, "parent_issue_id"),
			StartDate:      getStringPtrArg(args, "start_date"),
			DueDate:        getStringPtrArg(args, "due_date"),
			EstimatedHours: getFloat64PtrArg(args, "estimated_hours"),
			DoneRatio:      getIntPtrArg(args, "done_ratio"),
			Notes:          getStringArg(args, "notes"),
		}

		if err := client.UpdateIssue(ctx, id, params); err != nil {
			return errResult(err.Error()), nil
		}

		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated issue #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	}
}

func makeDeleteIssueHandler(client *api.Client) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return errResult("issue_id must be a positive integer"), nil
		}

		if err := client.DeleteIssue(ctx, id); err != nil {
			return errResult(err.Error()), nil
		}

		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted issue #%d", id)})
		if err != nil {
			return errResult(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return textResult(text), nil
	}
}
