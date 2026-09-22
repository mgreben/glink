package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":3000")
	t.Setenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/glink?sslmode=disable")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, ":3000", cfg.HTTP.Addr)
	assert.Equal(t, "postgres://user:password@localhost:5432/glink?sslmode=disable", cfg.Postgres.DSN)
	assert.Equal(t, "redis:6379", cfg.Redis.Addr)
	assert.Equal(t, time.Hour, cfg.Redis.LinkCacheTTL)
}

func TestNewConfigDefaultHTTPAddr(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/glink?sslmode=disable")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.HTTP.Addr)
}

func TestNewConfigLoadsRedisSettings(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/glink?sslmode=disable")
	t.Setenv("REDIS_ADDR", "localhost:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("LINK_CACHE_TTL", "30m")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, "localhost:6380", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 30*time.Minute, cfg.Redis.LinkCacheTTL)
}

func TestValidateMissingRequiredEnvVars(t *testing.T) {
	cfg := &Config{}

	err := cfg.validate()

	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "HTTP_ADDR"))
	assert.True(t, strings.Contains(err.Error(), "POSTGRES_DSN"))
}
