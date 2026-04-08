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

func NewCmdMcp(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server commands",
		Long:  "Model Context Protocol server for exposing Redmine operations to AI agents.",
	}

	cmd.AddCommand(newCmdServe(f))
	return cmd
}

func newCmdServe(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server",
		Long:  "Start a stdio-based MCP server that exposes Redmine operations as tools.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			s := server.NewMCPServer("rmn-redmine", "0.1.0")
			registerTools(s, client)

			return server.ServeStdio(s)
		},
	}

	return cmd
}

func registerTools(s *server.MCPServer, client *api.Client) {
	// list_issues
	s.AddTool(
		mcp.NewTool("list_issues",
			mcp.WithDescription("List Redmine issues with optional filters"),
			mcp.WithNumber("project_id", mcp.Description("Filter by project ID")),
			mcp.WithString("status_id", mcp.Description("Filter by status: open, closed, * (all), or a status ID")),
			mcp.WithString("assigned_to_id", mcp.Description("Filter by assignee: 'me' or a user ID")),
			mcp.WithNumber("tracker_id", mcp.Description("Filter by tracker ID")),
			mcp.WithNumber("limit", mcp.Description("Max number of results (default 25)")),
			mcp.WithNumber("offset", mcp.Description("Pagination offset")),
		),
		makeListIssuesHandler(client),
	)

	// get_issue
	s.AddTool(
		mcp.NewTool("get_issue",
			mcp.WithDescription("Get details of a specific Redmine issue"),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The issue ID")),
		),
		makeGetIssueHandler(client),
	)

	// create_issue
	s.AddTool(
		mcp.NewTool("create_issue",
			mcp.WithDescription("Create a new Redmine issue"),
			mcp.WithNumber("project_id", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("subject", mcp.Required(), mcp.Description("Issue subject")),
			mcp.WithString("description", mcp.Description("Issue description")),
			mcp.WithNumber("tracker_id", mcp.Description("Tracker ID")),
			mcp.WithNumber("priority_id", mcp.Description("Priority ID")),
			mcp.WithNumber("assigned_to_id", mcp.Description("Assignee user ID")),
		),
		makeCreateIssueHandler(client),
	)

	// update_issue
	s.AddTool(
		mcp.NewTool("update_issue",
			mcp.WithDescription("Update an existing Redmine issue"),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The issue ID to update")),
			mcp.WithString("subject", mcp.Description("New subject")),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithNumber("status_id", mcp.Description("New status ID")),
			mcp.WithNumber("priority_id", mcp.Description("New priority ID")),
			mcp.WithNumber("assigned_to_id", mcp.Description("New assignee user ID")),
			mcp.WithString("notes", mcp.Description("Add a note/comment to the issue")),
		),
		makeUpdateIssueHandler(client),
	)

	// delete_issue
	s.AddTool(
		mcp.NewTool("delete_issue",
			mcp.WithDescription("Delete a Redmine issue (cannot be undone)"),
			mcp.WithNumber("issue_id", mcp.Required(), mcp.Description("The issue ID to delete")),
		),
		makeDeleteIssueHandler(client),
	)
}

func toJSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
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

func makeListIssuesHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		params := api.IssueListParams{
			ProjectID:    getIntArg(args, "project_id"),
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

		return mcp.NewToolResultText(toJSON(result)), nil
	}
}

func makeGetIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id == 0 {
			return mcp.NewToolResultError("issue_id is required"), nil
		}

		issue, err := client.GetIssue(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(toJSON(issue)), nil
	}
}

func makeCreateIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		projectID := getIntArg(args, "project_id")
		subject := getStringArg(args, "subject")

		if projectID == 0 {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		if subject == "" {
			return mcp.NewToolResultError("subject is required"), nil
		}

		params := api.IssueCreateParams{
			ProjectID:    projectID,
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

		return mcp.NewToolResultText(toJSON(issue)), nil
	}
}

func makeUpdateIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id == 0 {
			return mcp.NewToolResultError("issue_id is required"), nil
		}

		params := api.IssueUpdateParams{
			Subject:      getStringArg(args, "subject"),
			Description:  getStringArg(args, "description"),
			StatusID:     getIntArg(args, "status_id"),
			PriorityID:   getIntArg(args, "priority_id"),
			AssignedToID: getIntArg(args, "assigned_to_id"),
			Notes:        getStringArg(args, "notes"),
		}

		if err := client.UpdateIssue(ctx, id, params); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(`{"status":"ok","message":"Updated issue #%d"}`, id)), nil
	}
}

func makeDeleteIssueHandler(client *api.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(req)

		id := getIntArg(args, "issue_id")
		if id == 0 {
			return mcp.NewToolResultError("issue_id is required"), nil
		}

		if err := client.DeleteIssue(ctx, id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(`{"status":"ok","message":"Deleted issue #%d"}`, id)), nil
	}
}
