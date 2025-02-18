package response

import (
	"book/web/templates"
	"context"
	"net/http"
	"strconv"
)

func NotFound(w http.ResponseWriter, msg string) {
	temlp := templates.Layout(templates.NotFound(msg), "NotFound")
	err := temlp.Render(context.Background(), w)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}

func InvalidInput(w http.ResponseWriter, msg string, statuesCode int) {
	temlp := templates.Layout(templates.InvalidInput(msg, strconv.Itoa(statuesCode)), "Invalid Input")
	err := temlp.Render(context.Background(), w)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}
