package project

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:     "delete <id-or-identifier>",
		Aliases: []string{"rm"},
		Short:   "Delete a project",
		Long:    "Delete a Redmine project. This action cannot be undone.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if !yes {
				fmt.Fprintf(f.IO.Out, "Delete project %s? This cannot be undone. [y/N]: ", id)
				reader := bufio.NewReader(f.IO.In)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("reading confirmation: %w", err)
				}
				if strings.ToLower(strings.TrimSpace(input)) != "y" {
					fmt.Fprintln(f.IO.Out, "Cancelled.")
					return nil
				}
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if err := client.DeleteProject(cmd.Context(), id); err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      string `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Deleted project %s", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Deleted project %s\n", id)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
