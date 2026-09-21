package app

import (
	"github.com/mgreben/glink/internal/config"
	"github.com/mgreben/glink/internal/handler"
	"github.com/mgreben/glink/internal/modules/links"
	"github.com/mgreben/glink/internal/modules/links/postgres"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			newDB,
			newValidator,
			postgres.NewLinkRepo,
			links.NewLinkService,
			handler.NewLinkHandler,
			handler.NewRouter,
			newServer,
		),
		fx.Invoke(registerHTTPServer),
	)
}
