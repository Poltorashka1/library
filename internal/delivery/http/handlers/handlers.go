package handlers

import (
	"book/internal/config"
	authorhandlers "book/internal/delivery/http/handlers/author"
	"book/internal/delivery/http/handlers/book"
	"book/internal/delivery/http/response"
	"book/internal/logger"
	"book/internal/usecase"
	"book/web/templates"
	"net/http"
)

// todo add logger middleware

//type Handlers interface {
//	bookhandlers.BookHandlers
//	authorhandlers.AuthorHandlers
//	NotFound(w http.ResponseWriter, r *http.Request)
//}

type Handlers struct {
	cfg config.HandlersConfig
	*BookHandlers
	*AuthorHandlers
}

func (h *Handlers) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	temlp := templates.Layout(templates.NotFound(""), "NotFound")
	err := temlp.Render(r.Context(), w)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func NewHandlers(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) *Handlers {
	return &Handlers{
		cfg,
		NewBookHandlers(log, cfg, useCase),
		NewAuthorHandlers(log, cfg, useCase),
	}
}

type BookHandlers struct {
	bookhandlers.GetBookHandler
	bookhandlers.CreateBookHandler
	bookhandlers.GetBooksHandler

	//ReadBook(w http.ResponseWriter, r *http.Request)
	//DownloadBook(w http.ResponseWriter, r *http.Request)
}

func NewBookHandlers(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) *BookHandlers {
	return &BookHandlers{
		bookhandlers.NewGetBookHandler(log, cfg, useCase),
		bookhandlers.NewCreateBookHandler(log, cfg, useCase),
		bookhandlers.NewGetBooksHandler(log, cfg, useCase),
	}
}

type AuthorHandlers struct {
	authorhandlers.GetAuthorHandler
}

func NewAuthorHandlers(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) *AuthorHandlers {
	return &AuthorHandlers{
		authorhandlers.NewGetAuthorHandler(log, cfg, useCase),
	}
}
