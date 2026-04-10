package wiki

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var projectID, text, comments string

	cmd := &cobra.Command{
		Use:     "create <title>",
		Aliases: []string{"new"},
		Short:   "Create a wiki page",
		Long:    "Create a new Redmine wiki page.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			if text == "" {
				return fmt.Errorf("--text is required")
			}
			title := args[0]

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.WikiPageCreateParams{
				Text:     text,
				Comments: comments,
			}

			page, err := client.CreateWikiPage(cmd.Context(), projectID, title, params)
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

			fmt.Fprintf(f.IO.Out, "Created wiki page %q (version %d)\n", page.Title, page.Version)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().StringVar(&text, "text", "", "Wiki page text (required)")
	cmd.Flags().StringVarP(&comments, "comments", "c", "", "Edit comments")

	return cmd
}
