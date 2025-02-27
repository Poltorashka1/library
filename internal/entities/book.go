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
	BookAuthors   BookAuthors
	BooksFileUUID sql.NullString
	CreatedAt     *time.Time
}

type BookAuthor struct {
	ID         sql.NullInt64
	UUID       sql.NullString
	NickName   sql.NullString
	Name       sql.NullString
	Surname    sql.NullString
	Patronymic sql.NullString
}

type BookAuthors []BookAuthor

type Books []Book

type BooksFilter struct {
	// Start first book id in the request
	Start int
	// Stop last book id in the request
	Stop int
	// BooksCount the number of books in the request
	BooksCount int
}

type BookFileFilter struct {
	FileName string
	FileType string
	Chapter  string
}

type BookFile struct {
	File []byte
}
