package auth

import (
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Manage Redmine authentication settings.",
	}

	cmd.AddCommand(NewCmdLogin(f))
	cmd.AddCommand(NewCmdStatus(f))

	return cmd
}
