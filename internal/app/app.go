package app

import (
	"github.com/mgreben/glink/internal/config"
	"github.com/mgreben/glink/internal/handler"
	"github.com/mgreben/glink/internal/modules/links"
	"github.com/mgreben/glink/internal/modules/links/postgres"
	redislinks "github.com/mgreben/glink/internal/modules/links/redis"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			newDB,
			newRedisClient,
			newValidator,
			postgres.NewLinkRepo,
			redislinks.NewLinkCache,
			links.NewLinkService,
			handler.NewLinkHandler,
			handler.NewRouter,
			newServer,
		),
		fx.Invoke(registerHTTPServer),
	)
}
