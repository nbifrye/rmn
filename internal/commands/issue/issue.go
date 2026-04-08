package issue

import (
	"encoding/json"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

// marshalJSON is used for JSON output. It can be replaced in tests.
var marshalJSON = json.MarshalIndent

func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Long:  "Work with Redmine issues.",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdClose(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
