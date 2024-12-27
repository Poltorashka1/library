package app

import (
	"book/internal/config"
	api "book/internal/delivery/http"
	"context"
)

type App struct {
	provider   *provider
	httpServer *api.Server
}

func New() *App {
	return &App{provider: NewProvider()}
}

func (a *App) Start(cfgFilename string) {
	config.MustLoad(cfgFilename)
	ctx := context.Background()

	a.httpServer = api.HTTPServer(a.provider.HttpConfig(), a.provider.Router(ctx), a.provider.Logger())
	a.httpServer.Start()
}
