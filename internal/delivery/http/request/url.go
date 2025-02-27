package request

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// URLParse parse url params
func URLParse(r *http.Request, key string) string {
	// input realisation there
	return chi.URLParam(r, key)
}
