package version

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var name, status, sharing, dueDate, description string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a version",
		Long:  "Update an existing Redmine project version.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid version ID: %s", args[0])
			}

			params := api.VersionUpdateParams{}
			changed := false

			if cmd.Flags().Changed("name") {
				params.Name = api.StringPtr(name)
				changed = true
			}
			if cmd.Flags().Changed("status") {
				params.Status = api.StringPtr(status)
				changed = true
			}
			if cmd.Flags().Changed("sharing") {
				params.Sharing = api.StringPtr(sharing)
				changed = true
			}
			if cmd.Flags().Changed("due-date") {
				params.DueDate = api.StringPtr(dueDate)
				changed = true
			}
			if cmd.Flags().Changed("description") {
				params.Description = api.StringPtr(description)
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

			if err := client.UpdateVersion(cmd.Context(), id, params); err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Updated version #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated version #%d\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Version name")
	cmd.Flags().StringVar(&status, "status", "", "Version status (open, locked, closed)")
	cmd.Flags().StringVar(&sharing, "sharing", "", "Version sharing (none, descendants, hierarchy, tree, system)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Version description")

	return cmd
}
