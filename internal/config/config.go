package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port   string `env:"PORT" envDefault:"8080"`
	DBPath string `env:"DB_PATH" envDefault:""`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
