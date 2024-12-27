package app

import (
	"book/internal/adapters/storage"
	"book/internal/config"
	api "book/internal/delivery/http"
	"book/internal/delivery/http/handlers"
	"book/internal/logger"
	"book/internal/repo"
	"book/internal/usecase"
	"context"
)

type provider struct {
	// todo in one struct
	httpCfg          config.HttpConfig
	dbConfig         config.DBConfig
	handlersConfig   config.HandlersConfig
	fileSystemConfig config.FileSystemConfig

	logger logger.Logger

	router api.Router

	handlers handlers.Handlers

	useCase usecase.UseCase

	repository repo.Repository

	db storage.DB
	fs storage.FS
}

func NewProvider() *provider {
	return &provider{}
}

func (p *provider) Logger() logger.Logger {
	if p.logger == nil {
		p.logger = logger.Load()
		p.Logger().Info("Logger loaded")
	}

	return p.logger
}

func (p *provider) Router(ctx context.Context) api.Router {
	if p.router == nil {
		p.router = api.NewHTTPRouter(p.Handlers(ctx))
		p.Logger().Info("Router loaded")
	}

	return p.router
}

func (p *provider) HttpConfig() config.HttpConfig {
	if p.httpCfg == nil {
		p.httpCfg = config.LoadHttpConfig()
		p.Logger().Info("HttpConfig loaded")
	}

	return p.httpCfg
}

func (p *provider) FileSystemConfig() config.FileSystemConfig {
	if p.fileSystemConfig == nil {
		p.fileSystemConfig = config.NewFileSystemConfig()
		p.Logger().Info("FileSystemConfig loaded")
	}

	return p.fileSystemConfig
}

func (p *provider) DBConfig() config.DBConfig {
	if p.dbConfig == nil {
		p.dbConfig = config.NewDbConfig()
		p.Logger().Info("DBConfig loaded")
	}

	return p.dbConfig
}

func (p *provider) Handlers(ctx context.Context) handlers.Handlers {
	if p.handlers == nil {
		p.handlers = handlers.NewHandlers(p.Logger(), p.UseCase(ctx), p.HandlersConfig())
		p.Logger().Info("Handlers loaded")
	}

	return p.handlers
}

func (p *provider) UseCase(ctx context.Context) usecase.UseCase {
	if p.useCase == nil {
		p.useCase = usecase.NewUseCase(p.Logger(), p.Repository(ctx))
		p.Logger().Info("UseCase loaded")
	}

	return p.useCase
}

func (p *provider) Repository(ctx context.Context) repo.Repository {
	if p.repository == nil {
		p.repository = repo.NewRepository(p.Logger(), p.DB(ctx), p.FS(ctx))
		p.Logger().Info("Repository loaded")
	}

	return p.repository
}

func (p *provider) DB(ctx context.Context) storage.DB {
	if p.db == nil {
		p.db = storage.NewDB(ctx, p.Logger(), p.DBConfig())
		p.Logger().Info("database loaded")
	}
	return p.db
}

func (p *provider) FS(ctx context.Context) storage.FS {
	if p.fs == nil {
		p.fs = storage.NewFileSystem(ctx, p.Logger(), p.FileSystemConfig())
		p.Logger().Info("file system loaded")
	}
	return p.fs
}
func (p *provider) HandlersConfig() config.HandlersConfig {
	if p.handlersConfig == nil {
		p.handlersConfig = config.NewHandlerConfig()
	}

	return p.handlersConfig
}
