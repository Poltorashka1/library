package bookhandlers

import (
	"book/internal/delivery/http/request"
	"book/internal/delivery/http/response"
	"net/http"
)

type Award struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

type Author struct {
	Name        string  `json:"name"`
	BirthYear   int     `json:"birthYear"`
	Nationality string  `json:"nationality"`
	Awards      []Award `json:"awards"`
}

type Ratings struct {
	Average float64 `json:"average"`
	Reviews int     `json:"reviews"`
}

type Chapter struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Published bool   `json:"published"`
}

type Book struct {
	Title         string `json:"title"`
	PublishedYear int    `json:"publishedYear"`
	Post          []string
	Genre         string    `json:"genre"`
	Author        Author    `json:"author"`
	Tests         []int     `json:"tests"`
	Ratings       Ratings   `json:"ratings"`
	Chapters      []Chapter `json:"chapters"`
}

func (h *bookHandlers) JSONTest(w http.ResponseWriter, r *http.Request) {
	payload := &Book{}

	err := request.JsonParseV2(r, payload)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
}
