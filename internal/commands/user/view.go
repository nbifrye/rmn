package user

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <id>",
		Aliases: []string{"show", "get"},
		Short:   "View a user",
		Long:    "Display details of a Redmine user. Use \"me\" to view the current user.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse and validate argument before making API calls.
			isMe := args[0] == "me"
			var id int
			if !isMe {
				var err error
				id, err = strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("invalid user ID: %s", args[0])
				}
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			var u *api.User

			if isMe {
				u, err = client.GetCurrentUser(cmd.Context())
				if err != nil {
					return err
				}
			} else {
				u, err = client.GetUser(cmd.Context(), id)
				if err != nil {
					return err
				}
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(u)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			admin := "No"
			if u.Admin {
				admin = "Yes"
			}

			fmt.Fprintf(f.IO.Out, "User #%d\n", u.ID)
			fmt.Fprintf(f.IO.Out, "Login:       %s\n", u.Login)
			fmt.Fprintf(f.IO.Out, "Name:        %s %s\n", u.FirstName, u.LastName)
			fmt.Fprintf(f.IO.Out, "Mail:        %s\n", u.Mail)
			fmt.Fprintf(f.IO.Out, "Admin:       %s\n", admin)
			if u.LastLoginOn != "" {
				fmt.Fprintf(f.IO.Out, "Last Login:  %s\n", u.LastLoginOn)
			}
			fmt.Fprintf(f.IO.Out, "Created:     %s\n", u.CreatedOn)
			return nil
		},
	}

	return cmd
}
