package handlers

import (
	"book/internal/config"
	"book/internal/delivery/http/handlers/author"
	"book/internal/delivery/http/handlers/book"
	"book/internal/delivery/http/response"
	"book/internal/logger"
	"book/internal/usecase"
	"book/web/templates"
	"net/http"
)

// todo add logger middleware

type Handlers interface {
	bookhandlers.BookHandlers
	authorhandlers.AuthorHandlers
	NotFound(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	config.HandlersConfig
	bookhandlers.BookHandlers
	authorhandlers.AuthorHandlers
}

func (h *handlers) NotFound(w http.ResponseWriter, r *http.Request) {
	nFTempl := templates.NotFound("")
	temlp := templates.Layout(nFTempl, "NotFound")
	err := temlp.Render(r.Context(), w)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func NewHandlers(logger logger.Logger, useCase usecase.UseCase, config config.HandlersConfig) Handlers {
	return &handlers{
		config,
		bookhandlers.NewBookHandlers(logger, useCase, config),
		authorhandlers.NewAuthorHandlers(logger, useCase, config),
	}
}
