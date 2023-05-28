package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB   DBConfig
	HTTP HTTP

	App     App
	Auth    Auth
	Metrics Metrics
}

func New() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process("ttto", c); err != nil {
		return nil, fmt.Errorf("envconfig proceed: %w", err)
	}

	return c, nil
}
