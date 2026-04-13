package project

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var name, identifier, description string
	var public bool
	var parentID int

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a project",
		Long:    "Create a new Redmine project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if identifier == "" {
				return fmt.Errorf("--identifier is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.ProjectCreateParams{
				Name:        name,
				Identifier:  identifier,
				Description: description,
				IsPublic:    public,
				ParentID:    parentID,
			}

			project, err := client.CreateProject(cmd.Context(), params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(project)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created project #%d: %s (%s)\n", project.ID, project.Name, project.Identifier)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Project name (required)")
	cmd.Flags().StringVar(&identifier, "identifier", "", "Project identifier (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Project description")
	cmd.Flags().BoolVar(&public, "public", false, "Make project public")
	cmd.Flags().IntVar(&parentID, "parent", 0, "Parent project ID")

	return cmd
}
