package membership

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
		Short:   "Delete a membership",
		Long:    "Delete a Redmine project membership. This action cannot be undone.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid membership ID: %s", args[0])
			}

			if !yes {
				fmt.Fprintf(f.IO.Out, "Delete membership #%d? This cannot be undone. [y/N]: ", id)
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

			if err := client.DeleteMembership(cmd.Context(), id); err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(struct {
					Status  string `json:"status"`
					ID      int    `json:"id"`
					Message string `json:"message"`
				}{Status: "ok", ID: id, Message: fmt.Sprintf("Deleted membership #%d", id)})
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Deleted membership #%d\n", id)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
