package project

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func projectStatusString(status int) string {
	switch status {
	case 1:
		return "active"
	case 5:
		return "archived"
	case 9:
		return "closed"
	default:
		return fmt.Sprintf("%d", status)
	}
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var status string
	var limit, offset int

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects",
		Long:    "List Redmine projects with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.ProjectListParams{
				Status: status,
				Limit:  limit,
				Offset: offset,
			}

			projects, total, err := client.ListProjects(cmd.Context(), params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Projects   []api.Project `json:"projects"`
					TotalCount int           `json:"total_count"`
				}{Projects: projects, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(projects) == 0 {
				fmt.Fprintln(f.IO.Out, "No projects found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tIDENTIFIER\tSTATUS\tPUBLIC")
			for _, p := range projects {
				publicStr := "no"
				if p.IsPublic {
					publicStr = "yes"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					p.ID,
					p.Name,
					p.Identifier,
					projectStatusString(p.Status),
					publicStr,
				)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d projects\n", len(projects), total)
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status (active, closed, archived)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Number of projects to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")

	return cmd
}
