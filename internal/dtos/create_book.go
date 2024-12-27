package dtos

type CreateBook struct {
	Name   string `json:"name"`
	Year   int    `json:"year"`
	Author string `json:"author"`
}
