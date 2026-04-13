package version

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
		Short:   "List versions",
		Long:    "List Redmine project versions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			versions, total, err := client.ListVersions(cmd.Context(), projectID)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Versions   []api.Version `json:"versions"`
					TotalCount int           `json:"total_count"`
				}{Versions: versions, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(versions) == 0 {
				fmt.Fprintln(f.IO.Out, "No versions found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDUE_DATE\tSHARING")
			for _, v := range versions {
				dueDate := ""
				if v.DueDate != nil {
					dueDate = *v.DueDate
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					v.ID,
					v.Name,
					v.Status,
					dueDate,
					v.Sharing,
				)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d versions\n", len(versions), total)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")

	return cmd
}
