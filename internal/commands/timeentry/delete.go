package timeentry

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a time entry",
		Long:    "Delete a Redmine time entry. This action cannot be undone.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid time entry ID: %s", args[0])
			}

			if !yes {
				fmt.Fprintf(f.IO.Out, "Delete time entry #%d? This cannot be undone. [y/N]: ", id)
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

			if err := client.DeleteTimeEntry(cmd.Context(), id); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Deleted time entry #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Deleted time entry #%d\n", id)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
