package issue

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var statusID, trackerID, priorityID, assignedToID, categoryID, versionID, parentID, doneRatio int
	var estimatedHours float64
	var subject, description, notes, startDate, dueDate string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an issue",
		Long:  "Update an existing Redmine issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			params := api.IssueUpdateParams{}
			changed := false

			if cmd.Flags().Changed("status") {
				params.StatusID = api.IntPtr(statusID)
				changed = true
			}
			if cmd.Flags().Changed("tracker") {
				params.TrackerID = api.IntPtr(trackerID)
				changed = true
			}
			if cmd.Flags().Changed("priority") {
				params.PriorityID = api.IntPtr(priorityID)
				changed = true
			}
			if cmd.Flags().Changed("subject") {
				params.Subject = api.StringPtr(subject)
				changed = true
			}
			if cmd.Flags().Changed("description") {
				params.Description = api.StringPtr(description)
				changed = true
			}
			if cmd.Flags().Changed("assignee") {
				params.AssignedToID = api.IntPtr(assignedToID)
				changed = true
			}
			if cmd.Flags().Changed("notes") {
				params.Notes = notes
				changed = true
			}
			if cmd.Flags().Changed("category") {
				params.CategoryID = api.IntPtr(categoryID)
				changed = true
			}
			if cmd.Flags().Changed("version") {
				params.FixedVersionID = api.IntPtr(versionID)
				changed = true
			}
			if cmd.Flags().Changed("parent") {
				params.ParentIssueID = api.IntPtr(parentID)
				changed = true
			}
			if cmd.Flags().Changed("start-date") {
				params.StartDate = api.StringPtr(startDate)
				changed = true
			}
			if cmd.Flags().Changed("due-date") {
				params.DueDate = api.StringPtr(dueDate)
				changed = true
			}
			if cmd.Flags().Changed("estimated-hours") {
				params.EstimatedHours = api.Float64Ptr(estimatedHours)
				changed = true
			}
			if cmd.Flags().Changed("done-ratio") {
				params.DoneRatio = api.IntPtr(doneRatio)
				changed = true
			}

			if !changed {
				fmt.Fprintln(f.IO.ErrOut, "No fields specified. Use flags to set fields to update.")
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.UpdateIssue(cmd.Context(), id, params); err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Updated issue #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated issue #%d\n", id)
			return nil
		},
	}

	cmd.Flags().IntVar(&statusID, "status", 0, "Status ID")
	cmd.Flags().IntVarP(&trackerID, "tracker", "t", 0, "Tracker ID")
	cmd.Flags().IntVar(&priorityID, "priority", 0, "Priority ID")
	cmd.Flags().StringVarP(&subject, "subject", "s", "", "Issue subject")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Issue description")
	cmd.Flags().IntVarP(&assignedToID, "assignee", "a", 0, "Assignee user ID")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Add a note/comment")
	cmd.Flags().IntVar(&categoryID, "category", 0, "Category ID")
	cmd.Flags().IntVar(&versionID, "version", 0, "Target version ID")
	cmd.Flags().IntVar(&parentID, "parent", 0, "Parent issue ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().Float64Var(&estimatedHours, "estimated-hours", 0, "Estimated hours")
	cmd.Flags().IntVar(&doneRatio, "done-ratio", 0, "Done ratio (0-100)")

	return cmd
}
