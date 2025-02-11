package bookhandlers

import (
	"book/internal/config"
	"book/internal/logger"
	"book/internal/usecase"
	"net/http"
)

//const (
//	JSON = "application/json"
//	PDF  = "application/pdf"
//	HTML = "text/html"
//)

type BookHandlers interface {
	Book(w http.ResponseWriter, r *http.Request)
	Books(w http.ResponseWriter, r *http.Request)
	CreateBook(w http.ResponseWriter, r *http.Request)
	UpdateBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
	ReadBook(w http.ResponseWriter, r *http.Request)
	DownloadBook(w http.ResponseWriter, r *http.Request)
	JSONTest(w http.ResponseWriter, r *http.Request)
}

type bookHandlers struct {
	logger  logger.Logger
	useCase usecase.UseCase
	cfg     config.HandlersConfig
}

func NewBookHandlers(logger logger.Logger, useCase usecase.UseCase, config config.HandlersConfig) BookHandlers {
	return &bookHandlers{logger: logger, useCase: useCase, cfg: config}
}
