package timeentry

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
		Short:   "View a time entry",
		Long:    "Display details of a Redmine time entry.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid time entry ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			entry, err := client.GetTimeEntry(cmd.Context(), id)
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

			fmt.Fprintf(f.IO.Out, "Time Entry #%d\n", entry.ID)
			fmt.Fprintf(f.IO.Out, "Project:     %s\n", entry.Project.Name)
			if entry.Issue != nil {
				fmt.Fprintf(f.IO.Out, "Issue:       #%d\n", entry.Issue.ID)
			}
			fmt.Fprintf(f.IO.Out, "User:        %s\n", entry.User.Name)
			fmt.Fprintf(f.IO.Out, "Activity:    %s\n", entry.Activity.Name)
			fmt.Fprintf(f.IO.Out, "Hours:       %.2f\n", entry.Hours)
			fmt.Fprintf(f.IO.Out, "Spent On:    %s\n", entry.SpentOn)
			if entry.Comments != "" {
				fmt.Fprintf(f.IO.Out, "Comments:    %s\n", entry.Comments)
			}
			fmt.Fprintf(f.IO.Out, "Created:     %s\n", entry.CreatedOn)
			fmt.Fprintf(f.IO.Out, "Updated:     %s\n", entry.UpdatedOn)
			return nil
		},
	}

	return cmd
}
