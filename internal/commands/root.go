package commands

import (
	"fmt"

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
	cmd.AddCommand(mcp.NewCmdMcp(f, version))
	cmd.AddCommand(newCmdCompletion())

	return cmd
}

func newCmdCompletion() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for the specified shell.

To load completions:

  bash:  source <(rmn completion bash)
  zsh:   rmn completion zsh > "${fpath[1]}/_rmn"
  fish:  rmn completion fish | source
  powershell: rmn completion powershell | Out-String | Invoke-Expression`,
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}
