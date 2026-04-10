package wiki

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var projectID string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List wiki pages",
		Long:    "List wiki pages in a Redmine project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			pages, err := client.ListWikiPages(cmd.Context(), projectID)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					WikiPages []api.WikiPage `json:"wiki_pages"`
				}{WikiPages: pages})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(pages) == 0 {
				fmt.Fprintln(f.IO.Out, "No wiki pages found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TITLE\tVERSION\tPARENT\tUPDATED")
			for _, p := range pages {
				parent := ""
				if p.Parent != nil {
					parent = p.Parent.Name
				}
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", p.Title, p.Version, parent, p.UpdatedOn)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")

	return cmd
}
