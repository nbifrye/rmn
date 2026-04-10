package user

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var status, limit, offset int
	var name string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List users",
		Long:    "List Redmine users with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.UserListParams{
				Status: status,
				Name:   name,
				Limit:  limit,
				Offset: offset,
			}

			users, total, err := client.ListUsers(cmd.Context(), params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Users      []api.User `json:"users"`
					TotalCount int        `json:"total_count"`
				}{Users: users, TotalCount: total})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(users) == 0 {
				fmt.Fprintln(f.IO.Out, "No users found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tLOGIN\tNAME\tMAIL\tADMIN")
			for _, u := range users {
				admin := "No"
				if u.Admin {
					admin = "Yes"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					u.ID,
					u.Login,
					u.FirstName+" "+u.LastName,
					u.Mail,
					admin,
				)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "\nShowing %d of %d users\n", len(users), total)
			return nil
		},
	}

	cmd.Flags().IntVar(&status, "status", 0, "Filter by status (0=active, 1=registered, 2=locked)")
	cmd.Flags().StringVar(&name, "name", "", "Filter by name or login")
	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Number of users to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")

	return cmd
}
