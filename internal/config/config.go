package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
}

type HTTPConfig struct {
	Addr string
}

type PostgresConfig struct {
	DSN string
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("HTTP_ADDR", ":8080")

	cfg := load(v)
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func load(v *viper.Viper) *Config {
	return &Config{
		HTTP: HTTPConfig{
			Addr: v.GetString("HTTP_ADDR"),
		},
		Postgres: PostgresConfig{
			DSN: v.GetString("POSTGRES_DSN"),
		},
	}
}

func (c *Config) validate() error {
	var missing []string

	if strings.TrimSpace(c.HTTP.Addr) == "" {
		missing = append(missing, "HTTP_ADDR")
	}
	if strings.TrimSpace(c.Postgres.DSN) == "" {
		missing = append(missing, "POSTGRES_DSN")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
