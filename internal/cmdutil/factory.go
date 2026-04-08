package cmdutil

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/config"
)

type Factory struct {
	Config    func() (*config.Config, error)
	APIClient func() (*api.Client, error)
	IO        *IOStreams

	flagURL    string
	flagAPIKey string
}

// SetFlagOverrides stores CLI flag values that override config file settings.
func (f *Factory) SetFlagOverrides(url, apiKey string) {
	f.flagURL = url
	f.flagAPIKey = apiKey
}

func NewFactory() *Factory {
	f := &Factory{
		IO: DefaultIOStreams(),
	}

	configFunc := func() (*config.Config, error) {
		return config.Load()
	}

	f.Config = configFunc
	f.APIClient = func() (*api.Client, error) {
		cfg, err := configFunc()
		if err != nil {
			return nil, err
		}
		if f.flagURL != "" {
			cfg.RedmineURL = f.flagURL
		}
		if f.flagAPIKey != "" {
			cfg.APIKey = f.flagAPIKey
		}
		if cfg.RedmineURL == "" || cfg.APIKey == "" {
			return nil, fmt.Errorf("not configured: run 'rmn auth login' to set up authentication")
		}
		return api.NewClient(cfg.RedmineURL, cfg.APIKey), nil
	}

	return f
}
