package project

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <id-or-identifier>",
		Aliases: []string{"show", "get"},
		Short:   "View a project",
		Long:    "Display details of a Redmine project.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			project, err := client.GetProject(cmd.Context(), id, nil)
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

			publicStr := "no"
			if project.IsPublic {
				publicStr = "yes"
			}

			fmt.Fprintf(f.IO.Out, "Project #%d\n", project.ID)
			fmt.Fprintf(f.IO.Out, "Name:        %s\n", project.Name)
			fmt.Fprintf(f.IO.Out, "Identifier:  %s\n", project.Identifier)
			fmt.Fprintf(f.IO.Out, "Status:      %s\n", projectStatusString(project.Status))
			fmt.Fprintf(f.IO.Out, "Public:      %s\n", publicStr)
			if project.Parent != nil {
				fmt.Fprintf(f.IO.Out, "Parent:      %s\n", project.Parent.Name)
			}
			if project.Homepage != "" {
				fmt.Fprintf(f.IO.Out, "Homepage:    %s\n", project.Homepage)
			}
			fmt.Fprintf(f.IO.Out, "Created:     %s\n", project.CreatedOn)
			fmt.Fprintf(f.IO.Out, "Updated:     %s\n", project.UpdatedOn)
			if project.Description != "" {
				fmt.Fprintf(f.IO.Out, "\n%s\n", project.Description)
			}
			return nil
		},
	}

	return cmd
}
