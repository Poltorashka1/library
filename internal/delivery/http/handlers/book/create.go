package bookhandlers

import (
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"book/internal/dtos"
	"errors"
	"fmt"
	"net/http"
)

func (h *bookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateBookRequest{}

	// todo Request Entity Too Large
	err := request.FormParse(r, payload)
	if err != nil {
		h.handleError(w, err)
		return
	}

	//err = payload.Validate()
	//if err != nil {
	//	response.Error(w, err, http.StatusUnprocessableEntity)
	//	return
	//}

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
	var multiErr *request.MultiError
	var fileNameErr *request.ErrFieldName
	var fileType *request.ErrInvalidFileType
	var maxSize *request.ErrContentToLarge
	var maxFieldSize *request.ErrFormValueToLarge
	switch {
	case errors.Is(err, request.ErrUnknownContentType):
		response.Error(w, err, http.StatusBadRequest)
		return
	case errors.As(err, &maxFieldSize):
		response.Error(w, err, http.StatusRequestEntityTooLarge)
		return
	case errors.As(err, &fileNameErr):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	case errors.As(err, &multiErr):
		response.Error(w, err, http.StatusBadRequest)
		return
	case errors.As(err, &fileType):
		response.Error(w, err, http.StatusUnsupportedMediaType)
		return
	case errors.As(err, &maxSize):
		response.Error(w, err, http.StatusRequestEntityTooLarge)
		return
	case errors.Is(err, request.ErrFileNameTooLong):
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	}
	h.logger.Error(err.Error())
	response.ServerError(w)
	return
}
