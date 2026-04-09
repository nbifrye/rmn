package issue

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <id>",
		Aliases: []string{"show", "get"},
		Short:   "View an issue",
		Long:  "Display details of a Redmine issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			issue, err := client.GetIssue(cmd.Context(), id, nil)
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

			assignee := "(unassigned)"
			if issue.AssignedTo != nil {
				assignee = issue.AssignedTo.Name
			}

			fmt.Fprintf(f.IO.Out, "Issue #%d\n", issue.ID)
			fmt.Fprintf(f.IO.Out, "Subject:     %s\n", issue.Subject)
			fmt.Fprintf(f.IO.Out, "Project:     %s\n", issue.Project.Name)
			fmt.Fprintf(f.IO.Out, "Tracker:     %s\n", issue.Tracker.Name)
			fmt.Fprintf(f.IO.Out, "Status:      %s\n", issue.Status.Name)
			fmt.Fprintf(f.IO.Out, "Priority:    %s\n", issue.Priority.Name)
			fmt.Fprintf(f.IO.Out, "Author:      %s\n", issue.Author.Name)
			fmt.Fprintf(f.IO.Out, "Assignee:    %s\n", assignee)
			if issue.Category != nil {
				fmt.Fprintf(f.IO.Out, "Category:    %s\n", issue.Category.Name)
			}
			if issue.FixedVersion != nil {
				fmt.Fprintf(f.IO.Out, "Version:     %s\n", issue.FixedVersion.Name)
			}
			if issue.Parent != nil {
				fmt.Fprintf(f.IO.Out, "Parent:      #%d\n", issue.Parent.ID)
			}
			if issue.StartDate != nil {
				fmt.Fprintf(f.IO.Out, "Start:       %s\n", *issue.StartDate)
			}
			if issue.DueDate != nil {
				fmt.Fprintf(f.IO.Out, "Due:         %s\n", *issue.DueDate)
			}
			if issue.EstimatedHours != nil {
				fmt.Fprintf(f.IO.Out, "Estimated:   %.2fh\n", *issue.EstimatedHours)
			}
			fmt.Fprintf(f.IO.Out, "Done:        %d%%\n", issue.DoneRatio)
			fmt.Fprintf(f.IO.Out, "Created:     %s\n", issue.CreatedOn)
			fmt.Fprintf(f.IO.Out, "Updated:     %s\n", issue.UpdatedOn)
			if issue.ClosedOn != nil {
				fmt.Fprintf(f.IO.Out, "Closed:      %s\n", *issue.ClosedOn)
			}
			if issue.Description != "" {
				fmt.Fprintf(f.IO.Out, "\n%s\n", issue.Description)
			}
			return nil
		},
	}

	return cmd
}
