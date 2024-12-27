package entities

import (
	"database/sql"
	"time"
)

type Book struct {
	ID              int
	UUID            string
	ISBN            string
	Title           string
	PublicationYear int
	Description     string
	//Image      string
	Publisher     string
	BookAuthors   Authors
	BooksFileUUID sql.NullString
	CreatedAt     *time.Time
}

type Books struct {
	Books []Book
}

type BookFilter struct {
	Start int
	Stop  int
}

type BookFileFilter struct {
	FileName string
	FileType string
	Chapter  string
}

type BookFile struct {
	File []byte
}
