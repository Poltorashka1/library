package authorhandlers

import (
	"book/internal/config"
	"book/internal/logger"
	"book/internal/usecase"
	authorusecase "book/internal/usecase/author"
	"net/http"
)

type GetAuthorHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type getAuthorHandler struct {
	log     logger.Logger
	cfg     config.HandlersConfig
	useCase authorusecase.GetAuthorUseCase
}

func NewGetAuthorHandler(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) GetAuthorHandler {
	return &getAuthorHandler{
		log, cfg, useCase.GetAuthorUseCase,
	}
}

func (h *getAuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}
