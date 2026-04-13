package membership

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <id>",
		Aliases: []string{"show", "get"},
		Short:   "View a membership",
		Long:    "Display details of a Redmine project membership.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid membership ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			m, err := client.GetMembership(cmd.Context(), id)
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(m)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Membership #%d\n", m.ID)
			fmt.Fprintf(f.IO.Out, "Project:     %s\n", m.Project.Name)
			if m.User != nil {
				fmt.Fprintf(f.IO.Out, "User:        %s (#%d)\n", m.User.Name, m.User.ID)
			}
			if m.Group != nil {
				fmt.Fprintf(f.IO.Out, "Group:       %s (#%d)\n", m.Group.Name, m.Group.ID)
			}
			roles := make([]string, 0, len(m.Roles))
			for _, r := range m.Roles {
				roles = append(roles, r.Name)
			}
			fmt.Fprintf(f.IO.Out, "Roles:       %s\n", strings.Join(roles, ", "))
			return nil
		},
	}

	return cmd
}
