package cmdutil

import (
	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/config"
)

type Factory struct {
	Config    func() (*config.Config, error)
	APIClient func() (*api.Client, error)
	IO        *IOStreams
}

func NewFactory() *Factory {
	io := DefaultIOStreams()

	configFunc := func() (*config.Config, error) {
		return config.Load()
	}

	return &Factory{
		Config: configFunc,
		APIClient: func() (*api.Client, error) {
			cfg, err := configFunc()
			if err != nil {
				return nil, err
			}
			return api.NewClient(cfg.RedmineURL, cfg.APIKey), nil
		},
		IO: io,
	}
}
