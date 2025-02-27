package bookhandlers

import (
	"book/internal/config"
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	apperrors "book/internal/errors"
	"book/internal/logger"
	"book/internal/usecase"
	bookusecase "book/internal/usecase/book"
	"book/web/templates"
	"errors"
	"net/http"
	"strconv"
)

type GetBooksHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type getBooksHandler struct {
	log     logger.Logger
	cfg     config.HandlersConfig
	useCase bookusecase.GetBooksUseCase
}

func NewGetBooksHandler(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) GetBooksHandler {
	return &getBooksHandler{
		log:     log,
		cfg:     cfg,
		useCase: useCase.GetBooksUseCase,
	}
}

func (h *getBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//todo error if page == 0 fix it
	// todo validate book limit if > max value
	var payload = &dtos.BooksRequest{
		Limit: h.cfg.DefaultBooksLimit(),
		Page:  h.cfg.DefaultBooksPageNumber(),
	}

	err := request.QueryParse(r, payload)
	if err != nil {
		// todo other error
		var mErr *request.MultiError
		switch {
		case errors.As(err, &mErr):
			response.InvalidInput(w, mErr.Error(), http.StatusUnprocessableEntity)
			return
		default:
			h.log.Error(err.Error())
			response.ServerError(w)
			return
		}
	}

	err = payload.Validate()
	if err != nil {
		response.InvalidInput(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// todo handle error
	books, err := h.useCase.Run(r.Context(), payload)
	if err != nil {
		switch {
		// todo refactor and check how it work
		case errors.Is(err, apperrors.ErrPageNotFound):
			response.Redirect(w, r, NewPath(r, "page", strconv.Itoa(h.cfg.DefaultBooksPageNumber())))
			return
		}
		h.log.Error(err.Error())
		response.ServerError(w)
		return
	}

	bookTempl := templates.Books(books, payload.Page, payload.Limit)
	temp := templates.Layout(bookTempl, "Books")
	err = temp.Render(r.Context(), w)
	if err != nil {
		// todo handle error
	}
}

func NewPath(r *http.Request, key string, val string) string {
	query := r.URL.Query()
	query.Set(key, val)
	r.URL.RawQuery = query.Encode()
	return r.URL.String()
}
