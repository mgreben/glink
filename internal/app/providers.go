package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mgreben/glink/internal/config"
	linksclickhouse "github.com/mgreben/glink/internal/modules/links/clickhouse"
	linkskafka "github.com/mgreben/glink/internal/modules/links/kafka"
	httpvalidator "github.com/mgreben/glink/pkg/http_validator"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func newDB(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func newRedisClient(lc fx.Lifecycle, cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return client
}

func newClickHouse(lc fx.Lifecycle, cfg *config.Config) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{Addr: []string{cfg.ClickHouse.Addr}})
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{OnStop: func(context.Context) error { return conn.Close() }})
	return conn, nil
}

func runClickConsumer(lc fx.Lifecycle, consumer *linkskafka.ClickConsumer, store *linksclickhouse.ClickStore) {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() { _ = consumer.Run(ctx, store) }()
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})
}

func newValidator() *validator.Validate {
	v := validator.New()

	_ = v.RegisterValidation("http_url", httpvalidator.ValidateHTTPURL)

	return v
}

func newServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
