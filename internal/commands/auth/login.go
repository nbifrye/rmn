package auth

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	var url, apiKey string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Redmine instance",
		Long:  "Configure the Redmine URL and API key for authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(f.IO.In)

			if url == "" {
				fmt.Fprint(f.IO.Out, "Redmine URL: ")
				input, _ := reader.ReadString('\n')
				url = strings.TrimSpace(input)
			}

			if apiKey == "" {
				fmt.Fprint(f.IO.Out, "API Key: ")
				input, _ := reader.ReadString('\n')
				apiKey = strings.TrimSpace(input)
			}

			if url == "" || apiKey == "" {
				return fmt.Errorf("both Redmine URL and API key are required")
			}

			cfg := &config.Config{
				RedmineURL: strings.TrimRight(url, "/"),
				APIKey:     apiKey,
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintln(f.IO.Out, "Authentication configured successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Redmine instance URL")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "Redmine API key")

	return cmd
}
