package membership

import (
	"encoding/json"
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

// marshalJSON marshals v to indented JSON. It is a variable so tests can replace it.
var marshalJSON = func(v interface{}) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling JSON: %w", err)
	}
	return data, nil
}

func NewCmdMembership(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "membership",
		Aliases: []string{"member"},
		Short:   "Manage project memberships",
		Long:    "Work with Redmine project memberships.",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
