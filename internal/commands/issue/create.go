package issue

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var trackerID, priorityID, assignedToID int
	var projectID, subject, description string

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create an issue",
		Long:  "Create a new Redmine issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			if subject == "" {
				return fmt.Errorf("--subject is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			var parsedProjectID interface{} = projectID
			if id, err := strconv.Atoi(projectID); err == nil {
				parsedProjectID = id
			}

			params := api.IssueCreateParams{
				ProjectID:    parsedProjectID,
				TrackerID:    trackerID,
				PriorityID:   priorityID,
				Subject:      subject,
				Description:  description,
				AssignedToID: assignedToID,
			}

			issue, err := client.CreateIssue(cmd.Context(), params)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, _ := json.MarshalIndent(issue, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created issue #%d: %s\n", issue.ID, issue.Subject)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().StringVarP(&subject, "subject", "s", "", "Issue subject (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Issue description")
	cmd.Flags().IntVarP(&trackerID, "tracker", "t", 0, "Tracker ID")
	cmd.Flags().IntVar(&priorityID, "priority", 0, "Priority ID")
	cmd.Flags().IntVarP(&assignedToID, "assignee", "a", 0, "Assignee user ID")

	return cmd
}
