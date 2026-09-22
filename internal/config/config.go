package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type HTTPConfig struct {
	Addr string
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr         string
	Password     string
	LinkCacheTTL time.Duration
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("REDIS_ADDR", "redis:6379")
	v.SetDefault("LINK_CACHE_TTL", "1h")

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
		Redis: RedisConfig{
			Addr:         v.GetString("REDIS_ADDR"),
			Password:     v.GetString("REDIS_PASSWORD"),
			LinkCacheTTL: v.GetDuration("LINK_CACHE_TTL"),
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
	if strings.TrimSpace(c.Redis.Addr) == "" {
		missing = append(missing, "REDIS_ADDR")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	if c.Redis.LinkCacheTTL <= 0 {
		return fmt.Errorf("LINK_CACHE_TTL must be greater than zero")
	}

	return nil
}
