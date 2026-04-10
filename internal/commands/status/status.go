package status

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

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Manage issue statuses",
		Long:  "Work with Redmine issue statuses.",
	}

	cmd.AddCommand(NewCmdList(f))

	return cmd
}
