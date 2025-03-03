package bookhandlers

import (
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	"book/internal/errors"
	"book/web/templates"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (h *bookHandlers) Book(w http.ResponseWriter, r *http.Request) {
	uuid := request.URLParse(r, "uuid")

	book, err := h.useCase.Book(r.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrBookNotFound):
			response.Redirect(w, r, h.cfg.NotFoundURL())
			return
		}
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	bookTempl := templates.Book(*book)
	templ := templates.Layout(bookTempl, book.Title)
	err = templ.Render(r.Context(), w)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	// response.Success(w, book, JSON)
}

func (h *bookHandlers) Books(w http.ResponseWriter, r *http.Request) {
	const op = "Books"
	//todo error if page == 0 fix it
	var payload = &dtos.BooksRequest{
		Limit: h.cfg.DefaultBooksLimit(),
		Page:  h.cfg.DefaultBooksPageNumber(),
	}
	fmt.Println(payload)

	err := request.QueryParse(r, payload)
	if err != nil {
		// todo other error
		var mErr *request.MultiError
		switch {
		case errors.As(err, &mErr):
			response.Error(w, mErr, http.StatusUnprocessableEntity)
			return
		default:
			h.logger.Error(err.Error())
			response.ServerError(w)
			return
		}
	}
	fmt.Printf("%+v\n", payload)
	// todo handle error
	books, err := h.useCase.Books(r.Context(), payload)
	if err != nil {
		switch {
		// todo refactor and check how it work
		case errors.Is(err, apperrors.ErrPageNotFound):
			response.Redirect(w, r, NewPath(r, "page", strconv.Itoa(h.cfg.DefaultBooksPageNumber())))
			return
		}
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	bookTempl := templates.Books(*books, payload.Page, payload.Limit)
	temp := templates.Layout(bookTempl, "Books")
	err = temp.Render(r.Context(), w)
	if err != nil {
		// todo handle error
	}
}

// todo подумать про указатели
func NewPath(r *http.Request, key string, val string) string {
	query := r.URL.Query()
	query.Set(key, val)
	r.URL.RawQuery = query.Encode()
	return r.URL.String()
}
