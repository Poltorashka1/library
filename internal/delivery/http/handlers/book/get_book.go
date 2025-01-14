package bookhandlers

import (
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	"book/internal/errors"
	"book/web/templates"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// todo delete
func parseQuery(query *url.URL, keys ...string) map[string]interface{} {
	queryList := make(map[string]interface{}, len(keys))
	re := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	for _, key := range keys {
		val := query.Query().Get(key)
		if re.MatchString(val) {
			strconv.Atoi(val)
		}
		queryList[key] = val
	}

	return queryList
}

func (h *bookHandlers) Book(w http.ResponseWriter, r *http.Request) {
	// todo chi.Parse
	// get book UUID from url

	parts := strings.Split(r.URL.Path, "/")
	uuid := parts[len(parts)-1]

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

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = h.cfg.BooksPageNumber()
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = h.cfg.BooksLimit()
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	var payload = dtos.BooksRequest{Limit: limit, Page: page}

	// todo handle error
	books, err := h.useCase.Books(r.Context(), payload)
	if err != nil {
		switch {
		// todo refactor and check how it work
		case errors.Is(err, apperrors.ErrPageNotFound):
			response.Redirect(w, r, NewPath(r, "page", h.cfg.BooksPageNumber()))
			return
		}
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	bookTempl := templates.Books(*books, page)
	temp := templates.Layout(bookTempl, "Books")
	err = temp.Render(r.Context(), w)
	if err != nil {
		// todo handle error
	}

	//tmpl, err := template.ParseFiles("book_template.html")
	//if err != nil {
	//	h.logger.ErrorOp(err.Error(), op)
	//	response.Error(w, err, http.StatusInternalServerError, JSON)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//err = tmpl.Execute(w, books)
	//if err != nil {
	//	h.logger.ErrorOp(err.Error(), op)
	//	response.Error(w, err, http.StatusInternalServerError, JSON)
	//	return
	//}

}

// todo подумать про указатели
func NewPath(r *http.Request, key string, val string) string {
	query := r.URL.Query()
	query.Set(key, val)
	r.URL.RawQuery = query.Encode()
	return r.URL.String()
}

//func (h *bookHandlers) Book2(w http.ResponseWriter, r *http.Request) {
//
//	parts := strings.Split(r.URL.Path, "/")
//	title := parts[len(parts)-1]
//	title = strings.ToLower(title)
//
//	book, err := h.useCase.BookInfo(r.Context(), title)
//	if err != nil {
//		response.Error(w, err, http.StatusBadRequest, JSON)
//		return
//	}
//	book.FilePath = filePath + book.FilePath
//
//	tmpl, err := template.ParseFiles("book_template.html")
//	if err != nil {
//		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	// Рендерим шаблон с данными о книге
//	err = tmpl.Execute(w, book)
//	if err != nil {
//		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
//		return
//	}
//}

//func (h *bookHandlers) Books(w http.ResponseWriter, r *http.Request) {
//	// todo
//}
