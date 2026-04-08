package commands

import (
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/commands/auth"
	"github.com/nbifrye/rmn/internal/commands/issue"
	"github.com/nbifrye/rmn/internal/commands/mcp"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rmn",
		Short:   "Redmine CLI",
		Long:    "rmn is a CLI tool for interacting with Redmine.",
		Version: version,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			url, _ := cmd.Root().PersistentFlags().GetString("redmine-url")
			key, _ := cmd.Root().PersistentFlags().GetString("api-key")
			if url != "" || key != "" {
				f.SetFlagOverrides(url, key)
			}
			return nil
		},
	}

	cmd.PersistentFlags().String("output", "table", "Output format: table or json")
	cmd.PersistentFlags().String("redmine-url", "", "Redmine instance URL (overrides config)")
	cmd.PersistentFlags().String("api-key", "", "Redmine API key (overrides config)")

	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(issue.NewCmdIssue(f))
	cmd.AddCommand(mcp.NewCmdMcp(f))

	return cmd
}
