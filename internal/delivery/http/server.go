package api

import (
	"book/internal/config"
	"book/internal/logger"
	"net/http"
)

type Server struct {
	cfg    config.HttpConfig
	router Router
	logger logger.Logger
}

func HTTPServer(cfg config.HttpConfig, router Router, logger logger.Logger) *Server {
	return &Server{cfg: cfg, router: router, logger: logger}
}

func (s *Server) Start() {
	s.logger.Info("Server started on " + s.cfg.Addr())
	err := http.ListenAndServe(s.cfg.Addr(), s.router)
	if err != nil {
		s.logger.Fatal(err.Error())
	}
}
