package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP
		Log
		PG
	}

	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	PG struct {
		DB          string `env-required:"true" env:"POSTGRES_DB"`
		User        string `env-required:"true" env:"POSTGRES_USER"`
		Password    string `env-required:"true" env:"POSTGRES_PASSWORD"`
		ServiceName string `env-required:"true" env:"POSTGRES_SERVICE_NAME"`
	}
)

// NewConfig returns app config.
func NewConfig() (cfg *Config, err error) {
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv error: %w", err)
	}

	return cfg, nil
}
