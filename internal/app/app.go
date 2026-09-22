package app

import (
	"github.com/mgreben/glink/internal/config"
	"github.com/mgreben/glink/internal/handler"
	"github.com/mgreben/glink/internal/modules/links"
	linksclickhouse "github.com/mgreben/glink/internal/modules/links/clickhouse"
	linkskafka "github.com/mgreben/glink/internal/modules/links/kafka"
	"github.com/mgreben/glink/internal/modules/links/postgres"
	redislinks "github.com/mgreben/glink/internal/modules/links/redis"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			newDB,
			newClickHouse,
			newRedisClient,
			newValidator,
			postgres.NewLinkRepo,
			linksclickhouse.NewClickStore,
			linkskafka.NewClickProducer,
			linkskafka.NewClickConsumer,
			redislinks.NewLinkCache,
			links.NewLinkServiceWithStats,
			handler.NewLinkHandler,
			handler.NewRouter,
			newServer,
		),
		fx.Invoke(registerHTTPServer, runClickConsumer),
	)
}
