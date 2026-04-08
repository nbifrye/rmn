package issue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an issue",
		Long:  "Delete a Redmine issue. This action cannot be undone.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.DeleteIssue(context.Background(), id); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Deleted issue #%d\n", id)
			return nil
		},
	}

	return cmd
}
