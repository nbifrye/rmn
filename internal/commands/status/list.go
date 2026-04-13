package status

import (
	"fmt"
	"text/tabwriter"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List issue statuses",
		Long:    "List all Redmine issue statuses.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			statuses, err := client.ListStatuses(cmd.Context())
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(statuses)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(statuses) == 0 {
				fmt.Fprintln(f.IO.Out, "No statuses found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tCLOSED")
			for _, s := range statuses {
				closed := "No"
				if s.IsClosed {
					closed = "Yes"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\n", s.ID, s.Name, closed)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			return nil
		},
	}

	return cmd
}
