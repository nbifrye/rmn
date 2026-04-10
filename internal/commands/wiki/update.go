package wiki

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var projectID, text, comments string
	var version int

	cmd := &cobra.Command{
		Use:   "update <title>",
		Short: "Update a wiki page",
		Long:  "Update an existing Redmine wiki page.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			title := args[0]

			params := api.WikiPageUpdateParams{}
			changed := false

			if cmd.Flags().Changed("text") {
				params.Text = api.StringPtr(text)
				changed = true
			}
			if cmd.Flags().Changed("comments") {
				params.Comments = api.StringPtr(comments)
				changed = true
			}
			if cmd.Flags().Changed("version") {
				params.Version = api.IntPtr(version)
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

			if err := client.UpdateWikiPage(cmd.Context(), projectID, title, params); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					Title   string `json:"title"`
					Message string `json:"message"`
				}{Status: "ok", Title: title, Message: fmt.Sprintf("Updated wiki page %q", title)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated wiki page %q\n", title)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().StringVar(&text, "text", "", "Wiki page text")
	cmd.Flags().StringVarP(&comments, "comments", "c", "", "Edit comments")
	cmd.Flags().IntVar(&version, "version", 0, "Expected current version (for optimistic locking)")

	return cmd
}
