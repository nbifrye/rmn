package timeentry

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var projectID, from, to string
	var issueID, userID, limit, offset int

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List time entries",
		Long:    "List Redmine time entries with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.TimeEntryListParams{
				ProjectID: projectID,
				IssueID:   issueID,
				UserID:    userID,
				From:      from,
				To:        to,
				Limit:     limit,
				Offset:    offset,
			}

			entries, total, err := client.ListTimeEntries(cmd.Context(), params)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					TimeEntries []api.TimeEntry `json:"time_entries"`
					TotalCount  int             `json:"total_count"`
				}{TimeEntries: entries, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(entries) == 0 {
				fmt.Fprintln(f.IO.Out, "No time entries found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tPROJECT\tISSUE\tUSER\tACTIVITY\tHOURS\tSPENT ON\tCOMMENTS")
			for _, e := range entries {
				issueStr := ""
				if e.Issue != nil {
					issueStr = fmt.Sprintf("#%d", e.Issue.ID)
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%.2f\t%s\t%s\n",
					e.ID, e.Project.Name, issueStr, e.User.Name, e.Activity.Name, e.Hours, e.SpentOn, e.Comments)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d time entries\n", len(entries), total)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Filter by project (ID or identifier)")
	cmd.Flags().IntVar(&issueID, "issue", 0, "Filter by issue ID")
	cmd.Flags().IntVar(&userID, "user", 0, "Filter by user ID")
	cmd.Flags().StringVar(&from, "from", "", "Filter from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&to, "to", "", "Filter to date (YYYY-MM-DD)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Number of entries to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")

	return cmd
}
