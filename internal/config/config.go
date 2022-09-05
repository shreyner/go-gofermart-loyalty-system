package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env   string `env:"ENV" envDefault:"dev"`
	Port  int    `env:"PORT" envDefault:"8080"`
	DBUrl string `env:"DATABASE_URL,required"`

	IsProd bool
	IsDev  bool
}

func (c *Config) Parse() error {
	if err := env.Parse(c); err != nil {
		return err
	}

	c.IsDev = c.Env == "dev"
	c.IsProd = c.Env != "dev"

	return nil
}
