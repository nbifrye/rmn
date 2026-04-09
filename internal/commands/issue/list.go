package issue

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var trackerID, limit, offset int
	var projectID, statusID, assignedToID, sort string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List issues",
		Long:  "List Redmine issues with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.IssueListParams{
				ProjectID:    projectID,
				StatusID:     statusID,
				AssignedToID: assignedToID,
				TrackerID:    trackerID,
				Sort:         sort,
				Limit:        limit,
				Offset:       offset,
			}

			issues, total, err := client.ListIssues(cmd.Context(), params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Issues     []api.Issue `json:"issues"`
					TotalCount int         `json:"total_count"`
				}{Issues: issues, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(issues) == 0 {
				fmt.Fprintln(f.IO.Out, "No issues found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTRACKER\tSTATUS\tPRIORITY\tASSIGNEE\tSUBJECT")
			for _, issue := range issues {
				assignee := ""
				if issue.AssignedTo != nil {
					assignee = issue.AssignedTo.Name
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					issue.ID,
					issue.Tracker.Name,
					issue.Status.Name,
					issue.Priority.Name,
					assignee,
					issue.Subject,
				)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d issues\n", len(issues), total)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Filter by project (ID or identifier)")
	cmd.Flags().StringVarP(&statusID, "status", "s", "", "Filter by status (open, closed, * for all, or status ID)")
	cmd.Flags().StringVarP(&assignedToID, "assignee", "a", "", "Filter by assignee (me or user ID)")
	cmd.Flags().IntVarP(&trackerID, "tracker", "t", 0, "Filter by tracker ID")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort by column (e.g. updated_on:desc, priority:asc)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Number of issues to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")

	return cmd
}
