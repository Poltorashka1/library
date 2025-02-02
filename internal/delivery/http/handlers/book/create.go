package bookhandlers

import (
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	"book/internal/errors"
	"errors"
	"fmt"
	"net/http"
)

func (h *bookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateBookRequest{}

	// todo add it in middleware
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.HTTPBodySize())
	if r.ContentLength/(1024*1024) > h.cfg.HTTPBodySize() {
		response.Error(w, &apperrors.ErrContentToLarge{MaxSize: h.cfg.HTTPBodySize()}, http.StatusRequestEntityTooLarge)
		return
	}

	// todo add error for ErrUnknownContentType
	err := request.FormParse(r, payload)
	if err != nil {
		if errors.Is(err, request.ErrUnknownContentType) {
			response.Error(w, err, http.StatusBadRequest)
			return
		}
		h.handleError(w, err)
		return
	}

	err = payload.Validate()
	if err != nil {
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("%+v\n", payload)

	result, err := h.useCase.CreateBook(r.Context(), payload)
	if err != nil {
		response.ServerError(w)
		h.logger.Error(err.Error())
		return
	}
	//os.Remove(payload.File.Name())
	response.Success(w, result, h.cfg.JSON())
}

func (h *bookHandlers) handleError(w http.ResponseWriter, err error) {
	mErr := &request.MultiError{}
	fErr := &request.ErrFieldType{}
	switch {
	case errors.As(err, &mErr):
		response.Error(w, err, http.StatusBadRequest)
		return
	case errors.As(err, &fErr):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	case errors.Is(err, request.ErrInvalidJsonSyntax):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	case errors.Is(err, request.ErrInvalidContentType):
		response.Error(w, err, http.StatusUnsupportedMediaType)
		return
	case errors.Is(err, request.ErrInvalidFieldType):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	case errors.Is(err, request.ErrContentToLarge):
		response.Error(w, &apperrors.ErrContentToLarge{MaxSize: h.cfg.HTTPBodySize()}, http.StatusRequestEntityTooLarge)
		return
	case errors.Is(err, request.ErrFileNameTooLong):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	case errors.Is(err, request.ErrFieldLength):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	}
	h.logger.Error(err.Error())
	response.ServerError(w)
	return
}
