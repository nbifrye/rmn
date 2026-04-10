package timeentry

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a time entry",
		Long:  "Update an existing Redmine time entry.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid time entry ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.TimeEntryUpdateParams{}

			if cmd.Flags().Changed("hours") {
				v, _ := cmd.Flags().GetFloat64("hours")
				params.Hours = api.Float64Ptr(v)
			}
			if cmd.Flags().Changed("activity") {
				v, _ := cmd.Flags().GetInt("activity")
				params.ActivityID = api.IntPtr(v)
			}
			if cmd.Flags().Changed("spent-on") {
				v, _ := cmd.Flags().GetString("spent-on")
				params.SpentOn = api.StringPtr(v)
			}
			if cmd.Flags().Changed("comments") {
				v, _ := cmd.Flags().GetString("comments")
				params.Comments = api.StringPtr(v)
			}

			if err := client.UpdateTimeEntry(cmd.Context(), id, params); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Updated time entry #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated time entry #%d\n", id)
			return nil
		},
	}

	cmd.Flags().Float64("hours", 0, "Hours spent")
	cmd.Flags().Int("activity", 0, "Activity ID")
	cmd.Flags().String("spent-on", "", "Date spent (YYYY-MM-DD)")
	cmd.Flags().StringP("comments", "c", "", "Comments")

	return cmd
}
