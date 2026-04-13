package tracker

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
		Short:   "List trackers",
		Long:    "List all Redmine trackers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			trackers, err := client.ListTrackers(cmd.Context())
			if err != nil {
				return err
			}

			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(trackers)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(trackers) == 0 {
				fmt.Fprintln(f.IO.Out, "No trackers found.")
				return nil
			}

			w := tabwriter.NewWriter(f.IO.Out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME")
			for _, t := range trackers {
				fmt.Fprintf(w, "%d\t%s\n", t.ID, t.Name)
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing output: %w", err)
			}

			return nil
		},
	}

	return cmd
}
