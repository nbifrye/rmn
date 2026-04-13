package project

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var name, description string
	var public bool
	var parentID int

	cmd := &cobra.Command{
		Use:   "update <id-or-identifier>",
		Short: "Update a project",
		Long:  "Update an existing Redmine project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			params := api.ProjectUpdateParams{}
			changed := false

			if cmd.Flags().Changed("name") {
				params.Name = api.StringPtr(name)
				changed = true
			}
			if cmd.Flags().Changed("description") {
				params.Description = api.StringPtr(description)
				changed = true
			}
			if cmd.Flags().Changed("public") {
				params.IsPublic = api.BoolPtr(public)
				changed = true
			}
			if cmd.Flags().Changed("parent") {
				params.ParentID = api.IntPtr(parentID)
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

			if err := client.UpdateProject(cmd.Context(), id, params); err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      string `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Updated project %s", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated project %s\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Project name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Project description")
	cmd.Flags().BoolVar(&public, "public", false, "Make project public")
	cmd.Flags().IntVar(&parentID, "parent", 0, "Parent project ID")

	return cmd
}
