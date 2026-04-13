package membership

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var projectID string
	var limit, offset int

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List project memberships",
		Long:    "List members and roles for a Redmine project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.MembershipListParams{Limit: limit, Offset: offset}
			memberships, total, err := client.ListMemberships(cmd.Context(), projectID, params)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Memberships []api.Membership `json:"memberships"`
					TotalCount  int              `json:"total_count"`
				}{Memberships: memberships, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(memberships) == 0 {
				fmt.Fprintln(f.IO.Out, "No memberships found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tUSER/GROUP\tROLES")
			for _, m := range memberships {
				name := "-"
				if m.User != nil {
					name = m.User.Name
				} else if m.Group != nil {
					name = "group: " + m.Group.Name
				}
				roles := make([]string, 0, len(m.Roles))
				for _, r := range m.Roles {
					roles = append(roles, r.Name)
				}
				fmt.Fprintf(w, "%d\t%s\t%s\n", m.ID, name, strings.Join(roles, ", "))
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d memberships\n", len(memberships), total)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Number of memberships to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")

	return cmd
}
