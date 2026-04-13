package timeentry

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var issueID, activityID int
	var hours float64
	var projectID, spentOn, comments string

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "log"},
		Short:   "Create a time entry",
		Long:    "Create a new Redmine time entry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if hours <= 0 {
				return fmt.Errorf("--hours is required and must be positive")
			}
			if issueID == 0 && projectID == "" {
				return fmt.Errorf("--issue or --project is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.TimeEntryCreateParams{
				IssueID:    issueID,
				ProjectID:  projectID,
				Hours:      hours,
				ActivityID: activityID,
				SpentOn:    spentOn,
				Comments:   comments,
			}

			entry, err := client.CreateTimeEntry(cmd.Context(), params)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(entry)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created time entry #%d: %.2fh\n", entry.ID, entry.Hours)
			return nil
		},
	}

	cmd.Flags().IntVar(&issueID, "issue", 0, "Issue ID")
	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier")
	cmd.Flags().Float64Var(&hours, "hours", 0, "Hours spent (required)")
	cmd.Flags().IntVar(&activityID, "activity", 0, "Activity ID")
	cmd.Flags().StringVar(&spentOn, "spent-on", "", "Date spent (YYYY-MM-DD, defaults to today)")
	cmd.Flags().StringVarP(&comments, "comments", "c", "", "Comments")

	return cmd
}
