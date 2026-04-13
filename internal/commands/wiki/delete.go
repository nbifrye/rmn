package wiki

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var projectID string
	var yes bool

	cmd := &cobra.Command{
		Use:     "delete <title>",
		Aliases: []string{"rm"},
		Short:   "Delete a wiki page",
		Long:    "Delete a Redmine wiki page. This action cannot be undone.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			title := args[0]

			if !yes {
				fmt.Fprintf(f.IO.Out, "Delete wiki page %q? This cannot be undone. [y/N]: ", title)
				reader := bufio.NewReader(f.IO.In)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("reading confirmation: %w", err)
				}
				if strings.ToLower(strings.TrimSpace(input)) != "y" {
					fmt.Fprintln(f.IO.Out, "Cancelled.")
					return nil
				}
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.DeleteWikiPage(cmd.Context(), projectID, title); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					Title   string `json:"title"`
					Message string `json:"message"`
				}{Status: "ok", Title: title, Message: fmt.Sprintf("Deleted wiki page %q", title)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Deleted wiki page %q\n", title)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
