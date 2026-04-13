package version

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <id>",
		Aliases: []string{"show", "get"},
		Short:   "View a version",
		Long:    "Display details of a Redmine project version.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid version ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			version, err := client.GetVersion(cmd.Context(), id)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(version)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			dueDate := ""
			if version.DueDate != nil {
				dueDate = *version.DueDate
			}

			fmt.Fprintf(f.IO.Out, "Version #%d\n", version.ID)
			fmt.Fprintf(f.IO.Out, "Name:        %s\n", version.Name)
			fmt.Fprintf(f.IO.Out, "Project:     %s\n", version.Project.Name)
			fmt.Fprintf(f.IO.Out, "Status:      %s\n", version.Status)
			fmt.Fprintf(f.IO.Out, "Sharing:     %s\n", version.Sharing)
			if dueDate != "" {
				fmt.Fprintf(f.IO.Out, "Due date:    %s\n", dueDate)
			}
			if version.WikiPageTitle != "" {
				fmt.Fprintf(f.IO.Out, "Wiki page:   %s\n", version.WikiPageTitle)
			}
			fmt.Fprintf(f.IO.Out, "Created:     %s\n", version.CreatedOn)
			fmt.Fprintf(f.IO.Out, "Updated:     %s\n", version.UpdatedOn)
			if version.Description != "" {
				fmt.Fprintf(f.IO.Out, "\n%s\n", version.Description)
			}
			return nil
		},
	}

	return cmd
}
