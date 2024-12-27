package authorhandlers

import (
	"book/internal/config"
	"book/internal/logger"
	"book/internal/usecase"
	"net/http"
)

type AuthorHandlers interface {
	Author(w http.ResponseWriter, r *http.Request)
	CreateAuthor(w http.ResponseWriter, r *http.Request)
	DeleteAuthor(w http.ResponseWriter, r *http.Request)
}

type authorHandlers struct {
	useCase usecase.UseCase
	logger  logger.Logger
	config  config.HandlersConfig
}

func (h *authorHandlers) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *authorHandlers) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func NewAuthorHandlers(logger logger.Logger, useCase usecase.UseCase, config config.HandlersConfig) AuthorHandlers {
	return &authorHandlers{
		useCase: useCase,
		logger:  logger,
		config:  config,
	}
}
