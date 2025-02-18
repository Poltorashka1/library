package bookhandlers

import (
	"book/internal/config"
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	"book/internal/errors"
	"book/internal/logger"
	"book/internal/usecase"
	"book/internal/usecase/book"
	"book/web/templates"
	"errors"
	"net/http"
)

type GetBookHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type getBookHandler struct {
	log     logger.Logger
	cfg     config.HandlersConfig
	useCase bookusecase.GetBookUseCase
}

func NewGetBookHandler(log logger.Logger, cfg config.HandlersConfig, useCase *usecase.UseCase) GetBookHandler {
	return &getBookHandler{
		log:     log,
		cfg:     cfg,
		useCase: useCase.GetBookUseCase,
	}
}

func (h *getBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid := request.URLParse(r, "uuid")
	err := h.uuidValidate(uuid)
	if errors.Is(err, apperrors.ErrInvalidBookID) {
		response.InvalidInput(w, "Invalid book id", http.StatusUnprocessableEntity)
		return
	}

	result, err := h.useCase.Run(r.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrBookNotFound):
			response.NotFound(w, "Book not found")
			return
		}
		response.ServerError(w)
		return
	}

	h.render(w, r, result)
}

func (h *getBookHandler) uuidValidate(uuid string) error {
	if len(uuid) != 36 {
		return apperrors.ErrInvalidBookID
	}

	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return apperrors.ErrInvalidBookID
	}

	for i, char := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if char != '-' {
				return apperrors.ErrInvalidBookID
			}
		} else {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return apperrors.ErrInvalidBookID
			}
		}
	}

	return nil
}

func (h *getBookHandler) render(w http.ResponseWriter, r *http.Request, result *dtos.BookResponse) {
	templ := templates.Layout(templates.Book(result), result.Title)
	err := templ.Render(r.Context(), w)
	if err != nil {
		h.log.Error(err.Error())
		response.ServerError(w)
		return
	}
}
