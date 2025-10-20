package config

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	Env          string        `env:"ENV" envDefault:"dev"`
}

func ParseConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
