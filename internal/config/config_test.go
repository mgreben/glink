package config

import (
	"strings"
	"testing"

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
}

func TestNewConfigDefaultHTTPAddr(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/glink?sslmode=disable")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.HTTP.Addr)
}

func TestValidateMissingRequiredEnvVars(t *testing.T) {
	cfg := &Config{}

	err := cfg.validate()

	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "HTTP_ADDR"))
	assert.True(t, strings.Contains(err.Error(), "POSTGRES_DSN"))
}
