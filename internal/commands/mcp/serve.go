package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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

// serveStdioFunc is the function used to start the MCP server. It can be replaced in tests.
var serveStdioFunc = func(s *server.MCPServer) error {
	return server.ServeStdio(s)
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

			s := server.NewMCPServer("rmn-redmine", version)
			registerTools(s, client)

			return serveStdioFunc(s)
		},
	}

	return cmd
}

func registerTools(s *server.MCPServer, client *api.Client) {
	// list_issues — read-only, idempotent
	s.AddTool(
		mcp.NewTool("list_issues",
			mcp.WithDescription("List Redmine issues with optional filters. Returns a JSON object with an 'issues' array and 'total_count'. Without filters, returns open issues. Use 'offset' and 'limit' for pagination through large result sets."),
			mcp.WithTitleAnnotation("List Redmine Issues"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithString("project_id", mcp.Description("Filter by project ID (numeric) or identifier (string, e.g. 'my-project')")),
			mcp.WithString("status_id", mcp.Description("Filter by status: 'open' (default), 'closed', '*' (all), or a numeric status ID")),
			mcp.WithString("assigned_to_id", mcp.Description("Filter by assignee: 'me' for current user, or a numeric user ID")),
			mcp.WithNumber("tracker_id", mcp.Description("Filter by tracker ID (values are specific to your Redmine instance)")),
			mcp.WithNumber("limit", mcp.Description("Max number of results to return (default 25, max 100)")),
			mcp.WithNumber("offset", mcp.Description("Pagination offset for retrieving subsequent pages")),
		),
		makeListIssuesHandler(client),
	)

	// get_issue — read-only, idempotent
	s.AddTool(
		mcp.NewTool("get_issue",
			mcp.WithDescription("Get full details of a specific Redmine issue by ID, including project, tracker, status, priority, author, assignee, description, and timestamps."),
			mcp.WithTitleAnnotation("Get Redmine Issue"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The numeric issue ID")),
		),
		makeGetIssueHandler(client),
	)

	// create_issue — not read-only, not idempotent (creates new resource each time)
	s.AddTool(
		mcp.NewTool("create_issue",
			mcp.WithDescription("Create a new Redmine issue. Returns the created issue with its assigned ID."),
			mcp.WithTitleAnnotation("Create Redmine Issue"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID (numeric) or identifier (string, e.g. 'my-project')")),
			mcp.WithString("subject", mcp.Required(), mcp.Description("Issue subject/title")),
			mcp.WithString("description", mcp.Description("Detailed issue description (supports Textile markup)")),
			mcp.WithNumber("tracker_id", mcp.Description("Tracker ID (values are specific to your Redmine instance)")),
			mcp.WithNumber("priority_id", mcp.Description("Priority ID (values are specific to your Redmine instance)")),
			mcp.WithNumber("assigned_to_id", mcp.Description("User ID to assign the issue to")),
		),
		makeCreateIssueHandler(client),
	)

	// update_issue — not read-only, idempotent (same update yields same result)
	s.AddTool(
		mcp.NewTool("update_issue",
			mcp.WithDescription("Update an existing Redmine issue. Only provided fields are changed; omitted fields are left unchanged."),
			mcp.WithTitleAnnotation("Update Redmine Issue"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The numeric issue ID to update")),
			mcp.WithString("subject", mcp.Description("New subject/title")),
			mcp.WithString("description", mcp.Description("New description (supports Textile markup)")),
			mcp.WithNumber("status_id", mcp.Description("New status ID (values are specific to your Redmine instance)")),
			mcp.WithNumber("priority_id", mcp.Description("New priority ID (values are specific to your Redmine instance)")),
			mcp.WithNumber("assigned_to_id", mcp.Description("New assignee user ID (set to 0 to unassign)")),
			mcp.WithString("notes", mcp.Description("Add a comment/note to the issue")),
		),
		makeUpdateIssueHandler(client),
	)

	// delete_issue — destructive, idempotent
	s.AddTool(
		mcp.NewTool("delete_issue",
			mcp.WithDescription("Permanently delete a Redmine issue. This action cannot be undone. No confirmation is requested."),
			mcp.WithTitleAnnotation("Delete Redmine Issue"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The numeric issue ID to delete")),
		),
		makeDeleteIssueHandler(client),
	)
}

func toJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(data), nil
}

func getArgs(req mcp.CallToolRequest) map[string]interface{} {
	if m, ok := req.Params.Arguments.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
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

func makeListIssuesHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		params := api.IssueListParams{
			ProjectID:    getStringArg(args, "project_id"),
			StatusID:     getStringArg(args, "status_id"),
			AssignedToID: getStringArg(args, "assigned_to_id"),
			TrackerID:    getIntArg(args, "tracker_id"),
			Limit:        getIntArg(args, "limit"),
			Offset:       getIntArg(args, "offset"),
		}

		issues, total, err := client.ListIssues(ctx, params)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := struct {
			Issues     []api.Issue `json:"issues"`
			TotalCount int         `json:"total_count"`
		}{Issues: issues, TotalCount: total}

		text, err := toJSON(result)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return mcp.NewToolResultText(text), nil
	}
}

func makeGetIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return mcp.NewToolResultError("issue_id must be a positive integer"), nil
		}

		issue, err := client.GetIssue(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		text, err := toJSON(issue)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return mcp.NewToolResultText(text), nil
	}
}

func makeCreateIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		projectID := getStringArg(args, "project_id")
		subject := getStringArg(args, "subject")

		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		if subject == "" {
			return mcp.NewToolResultError("subject is required"), nil
		}

		var parsedProjectID interface{} = projectID
		if id, err := strconv.Atoi(projectID); err == nil {
			parsedProjectID = id
		}

		params := api.IssueCreateParams{
			ProjectID:    parsedProjectID,
			Subject:      subject,
			Description:  getStringArg(args, "description"),
			TrackerID:    getIntArg(args, "tracker_id"),
			PriorityID:   getIntArg(args, "priority_id"),
			AssignedToID: getIntArg(args, "assigned_to_id"),
		}

		issue, err := client.CreateIssue(ctx, params)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		text, err := toJSON(issue)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return mcp.NewToolResultText(text), nil
	}
}

func makeUpdateIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return mcp.NewToolResultError("issue_id must be a positive integer"), nil
		}

		params := api.IssueUpdateParams{
			Subject:      getStringArg(args, "subject"),
			Notes:        getStringArg(args, "notes"),
			Description:  getStringPtrArg(args, "description"),
			StatusID:     getIntPtrArg(args, "status_id"),
			PriorityID:   getIntPtrArg(args, "priority_id"),
			AssignedToID: getIntPtrArg(args, "assigned_to_id"),
		}

		if err := client.UpdateIssue(ctx, id, params); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Updated issue #%d", id)})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return mcp.NewToolResultText(text), nil
	}
}

func makeDeleteIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id <= 0 {
			return mcp.NewToolResultError("issue_id must be a positive integer"), nil
		}

		if err := client.DeleteIssue(ctx, id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		text, err := toJSON(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{Status: "ok", Message: fmt.Sprintf("Deleted issue #%d", id)})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}
		return mcp.NewToolResultText(text), nil
	}
}
