package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" required:"true" default:"account-service"`
	AppEnv      string `env:"APP_ENV" required:"true" default:"development"`
	Host        string `env:"HTTP_HOST" required:"true" default:"localhost"`
	Port        int    `env:"HTTP_PORT" required:"true" default:"9000"`
	LogLevel    string `env:"LOG_LEVEL" required:"true" default:"info"`
	DbDsn       string `env:"DB_DSN" required:"true"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
