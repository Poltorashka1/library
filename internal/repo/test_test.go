package repo

import (
	"book/internal/adapters/storage"
	"book/internal/config"
	"book/internal/logger"
	"context"
	"testing"
)

type setup struct {
	cfg    config.DBConfig
	db     storage.DB
	logger logger.Logger
	r      Repository
	ctx    context.Context
}

var s = setup{}

func TestMain(m *testing.M) {
	config.MustLoad("C:\\Users\\Sanfo\\GolandProjects\\Arch\\.env")

	s.cfg = config.NewDbConfig()
	s.logger = logger.Load()
	s.db = storage.NewDB(context.Background(), s.logger, s.cfg)
	s.r = NewRepository(s.logger, s.db, nil)
	s.ctx = context.Background()
	m.Run()
}
