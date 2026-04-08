package issue

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var statusID, trackerID, priorityID, assignedToID int
	var subject, description, notes string

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

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.IssueUpdateParams{
				StatusID:     statusID,
				TrackerID:    trackerID,
				PriorityID:   priorityID,
				Subject:      subject,
				Description:  description,
				AssignedToID: assignedToID,
				Notes:        notes,
			}

			if err := client.UpdateIssue(cmd.Context(), id, params); err != nil {
				return err
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

	return cmd
}
