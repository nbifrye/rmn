package auth

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  "Display the current Redmine authentication configuration and verify connectivity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			if cfg.RedmineURL == "" {
				fmt.Fprintln(f.IO.Out, "Not configured. Run 'rmn auth login' to set up authentication.")
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Redmine URL: %s\n", cfg.RedmineURL)
			if len(cfg.APIKey) > 4 {
				fmt.Fprintf(f.IO.Out, "API Key:     %s***\n", cfg.APIKey[:4])
			} else if cfg.APIKey != "" {
				fmt.Fprintln(f.IO.Out, "API Key:     ***")
			} else {
				fmt.Fprintln(f.IO.Out, "API Key:     (not set)")
			}

			if cfg.APIKey == "" {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				fmt.Fprintf(f.IO.ErrOut, "Connection failed: %v\n", err)
				return nil
			}

			var result struct {
				User struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
				} `json:"user"`
			}

			err = client.Get(cmd.Context(), "/users/current.json", nil, &result)
			if err != nil {
				fmt.Fprintf(f.IO.ErrOut, "Connection failed: %v\n", err)
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Logged in as: %s (ID: %d)\n", result.User.Login, result.User.ID)
			return nil
		},
	}

	return cmd
}
