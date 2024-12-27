package authorhandlers

import "net/http"

func (h *authorHandlers) Author(w http.ResponseWriter, r *http.Request) {
	h.useCase.Author(123)
}
