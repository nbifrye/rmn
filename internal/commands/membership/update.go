package membership

import (
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var roleIDs []int

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a membership's roles",
		Long:  "Update the role IDs of an existing Redmine project membership.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid membership ID: %s", args[0])
			}
			if len(roleIDs) == 0 {
				return fmt.Errorf("--role is required (at least one)")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.MembershipUpdateParams{RoleIDs: roleIDs}
			if err := client.UpdateMembership(cmd.Context(), id, params); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Updated membership #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Updated membership #%d\n", id)
			return nil
		},
	}

	cmd.Flags().IntSliceVar(&roleIDs, "role", nil, "Role ID (required, repeatable)")

	return cmd
}
