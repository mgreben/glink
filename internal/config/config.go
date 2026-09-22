package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP       HTTPConfig
	Postgres   PostgresConfig
	Redis      RedisConfig
	Kafka      KafkaConfig
	ClickHouse ClickHouseConfig
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

type KafkaConfig struct {
	Brokers            []string
	ClickTopic         string
	ClickConsumerGroup string
	MaxBufferedRecords int
}

type ClickHouseConfig struct {
	Addr string
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("REDIS_ADDR", "redis:6379")
	v.SetDefault("LINK_CACHE_TTL", "1h")
	v.SetDefault("KAFKA_BROKERS", "kafka:19092")
	v.SetDefault("KAFKA_CLICK_TOPIC", "link.clicked")
	v.SetDefault("KAFKA_CLICK_CONSUMER_GROUP", "glink-stats-v1")
	v.SetDefault("KAFKA_MAX_BUFFERED_RECORDS", 10000)
	v.SetDefault("CLICKHOUSE_ADDR", "clickhouse:9000")

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
		Kafka: KafkaConfig{
			Brokers:            strings.Split(v.GetString("KAFKA_BROKERS"), ","),
			ClickTopic:         v.GetString("KAFKA_CLICK_TOPIC"),
			ClickConsumerGroup: v.GetString("KAFKA_CLICK_CONSUMER_GROUP"),
			MaxBufferedRecords: v.GetInt("KAFKA_MAX_BUFFERED_RECORDS"),
		},
		ClickHouse: ClickHouseConfig{Addr: v.GetString("CLICKHOUSE_ADDR")},
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
	if len(c.Kafka.Brokers) == 0 || strings.TrimSpace(c.Kafka.Brokers[0]) == "" {
		missing = append(missing, "KAFKA_BROKERS")
	}
	if strings.TrimSpace(c.Kafka.ClickTopic) == "" {
		missing = append(missing, "KAFKA_CLICK_TOPIC")
	}
	if strings.TrimSpace(c.Kafka.ClickConsumerGroup) == "" {
		missing = append(missing, "KAFKA_CLICK_CONSUMER_GROUP")
	}
	if strings.TrimSpace(c.ClickHouse.Addr) == "" {
		missing = append(missing, "CLICKHOUSE_ADDR")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	if c.Redis.LinkCacheTTL <= 0 {
		return fmt.Errorf("LINK_CACHE_TTL must be greater than zero")
	}
	if c.Kafka.MaxBufferedRecords <= 0 {
		return fmt.Errorf("KAFKA_MAX_BUFFERED_RECORDS must be greater than zero")
	}

	return nil
}
