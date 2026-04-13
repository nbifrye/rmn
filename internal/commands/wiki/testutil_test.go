package wiki

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/nbifrye/rmn/internal/config"
	"github.com/spf13/cobra"
)

func newTestFactory(srv *httptest.Server) *cmdutil.Factory {
	return &cmdutil.Factory{
		Config: func() (*config.Config, error) {
			return &config.Config{RedmineURL: srv.URL, APIKey: "test"}, nil
		},
		APIClient: func() (*api.Client, error) {
			return api.NewClient(srv.URL, "test"), nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}
}

// newNoServerFactory creates a Factory for tests where no HTTP call should be made.
func newNoServerFactory(t *testing.T) *cmdutil.Factory {
	t.Helper()
	return &cmdutil.Factory{
		Config: func() (*config.Config, error) { return &config.Config{}, nil },
		APIClient: func() (*api.Client, error) {
			t.Fatal("APIClient should not be called in this test")
			return nil, nil
		},
		IO: &cmdutil.IOStreams{
			In:     &bytes.Buffer{},
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}
}

func setupRootFlags(cmd *cobra.Command, output string) {
	cmd.Root().PersistentFlags().String("output", output, "")
	cmd.Root().PersistentFlags().String("redmine-url", "", "")
	cmd.Root().PersistentFlags().String("api-key", "", "")
}
