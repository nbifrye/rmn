package wiki

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	var projectID string
	var version int

	cmd := &cobra.Command{
		Use:     "view <title>",
		Aliases: []string{"show", "get"},
		Short:   "View a wiki page",
		Long:    "Display the content of a Redmine wiki page.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			title := args[0]

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			page, err := client.GetWikiPage(cmd.Context(), projectID, title, version)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(page)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Title:    %s\n", page.Title)
			fmt.Fprintf(f.IO.Out, "Version:  %d\n", page.Version)
			if page.Parent != nil {
				fmt.Fprintf(f.IO.Out, "Parent:   %s\n", page.Parent.Name)
			}
			fmt.Fprintf(f.IO.Out, "Author:   %s\n", page.Author.Name)
			fmt.Fprintf(f.IO.Out, "Created:  %s\n", page.CreatedOn)
			fmt.Fprintf(f.IO.Out, "Updated:  %s\n", page.UpdatedOn)
			if page.Comments != "" {
				fmt.Fprintf(f.IO.Out, "Comments: %s\n", page.Comments)
			}
			fmt.Fprintf(f.IO.Out, "\n%s\n", page.Text)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().IntVar(&version, "version", 0, "Specific version to retrieve")

	return cmd
}
