package tracker

import (
	"encoding/json"
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

var marshalJSON = func(v interface{}) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling JSON: %w", err)
	}
	return data, nil
}

func NewCmdTracker(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tracker",
		Short: "Manage trackers",
		Long:  "Work with Redmine trackers.",
	}

	cmd.AddCommand(NewCmdList(f))

	return cmd
}
