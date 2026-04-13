package membership

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var projectID string
	var userID int
	var roleIDs []int

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "add"},
		Short:   "Add a member to a project",
		Long:    "Create a new Redmine project membership.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			if userID == 0 {
				return fmt.Errorf("--user is required")
			}
			if len(roleIDs) == 0 {
				return fmt.Errorf("--role is required (at least one)")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.MembershipCreateParams{
				UserID:  userID,
				RoleIDs: roleIDs,
			}

			m, err := client.CreateMembership(cmd.Context(), projectID, params)
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

			fmt.Fprintf(f.IO.Out, "Created membership #%d\n", m.ID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().IntVar(&userID, "user", 0, "User ID (required)")
	cmd.Flags().IntSliceVar(&roleIDs, "role", nil, "Role ID (required, repeatable)")

	return cmd
}
