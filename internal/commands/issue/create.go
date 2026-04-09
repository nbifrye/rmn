package issue

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var trackerID, priorityID, assignedToID, categoryID, versionID, parentID, doneRatio int
	var estimatedHours float64
	var projectID, subject, description, startDate, dueDate string

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
				ProjectID:      parsedProjectID,
				TrackerID:      trackerID,
				PriorityID:     priorityID,
				Subject:        subject,
				Description:    description,
				AssignedToID:   assignedToID,
				CategoryID:     categoryID,
				FixedVersionID: versionID,
				ParentIssueID:  parentID,
				StartDate:      startDate,
				DueDate:        dueDate,
				EstimatedHours: estimatedHours,
				DoneRatio:      doneRatio,
			}

			issue, err := client.CreateIssue(cmd.Context(), params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(issue)
				if err != nil {
					return err
				}
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
	cmd.Flags().IntVar(&categoryID, "category", 0, "Category ID")
	cmd.Flags().IntVar(&versionID, "version", 0, "Target version ID")
	cmd.Flags().IntVar(&parentID, "parent", 0, "Parent issue ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().Float64Var(&estimatedHours, "estimated-hours", 0, "Estimated hours")
	cmd.Flags().IntVar(&doneRatio, "done-ratio", 0, "Done ratio (0-100)")

	return cmd
}
