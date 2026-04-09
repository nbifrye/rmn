package issue

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

const defaultClosedStatusID = 5

func NewCmdClose(f *cmdutil.Factory) *cobra.Command {
	var statusID int
	var notes string

	cmd := &cobra.Command{
		Use:   "close <id>",
		Short: "Close an issue",
		Long:  "Close a Redmine issue by setting its status to closed. Uses status ID 5 (Closed) by default.",
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
				StatusID: api.IntPtr(statusID),
			}
			if cmd.Flags().Changed("notes") {
				params.Notes = notes
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
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Closed issue #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Closed issue #%d\n", id)
			return nil
		},
	}

	cmd.Flags().IntVar(&statusID, "status", defaultClosedStatusID, "Closed status ID")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Add a closing note")

	return cmd
}
