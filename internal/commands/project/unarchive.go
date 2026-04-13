package project

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUnarchive(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unarchive <id-or-identifier>",
		Short: "Unarchive a project",
		Long:  "Unarchive a Redmine project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.UnarchiveProject(cmd.Context(), id); err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      string `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Unarchived project %s", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Unarchived project %s\n", id)
			return nil
		},
	}

	return cmd
}
