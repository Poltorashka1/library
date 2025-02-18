package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	apperrors "book/internal/errors"
	"context"
	"strings"
)

// Books get books with authors, genres; err return apperrors.ErrPageNotFound
func (r *bookRepository) Books(ctx context.Context, filter entities.BookFilter) (*entities.Books, error) {
	// todo mb in cache

	query := storage.Query{QueryName: "get books count", Query: "select count(*) from books"}
	var booksCount int
	err := r.db.QueryRowContext(ctx, query).Scan(&booksCount)
	if err != nil {
		return nil, err
	}
	if filter.Start > booksCount {
		return nil, apperrors.ErrPageNotFound
	}

	query = storage.Query{
		QueryName: "select books with pagination",
		Query: `
				SELECT t1.id,t1.uuid,t1.isbn,t1.title,t1.publication_year,t1.description,
       				GROUP_CONCAT(t3.nickname, ', ') AS authors,
       				GROUP_CONCAT(t3.uuid, ', ')     AS authors_uuid
				FROM books AS t1
         			JOIN book_author AS t2 ON t2.book_id = t1.id
         			JOIN authors AS t3 ON t2.author_id = t3.id
				WHERE t1.id > $1 and t1.id <= $2
				GROUP BY t1.id
				ORDER BY t1.id`,
		Args: []any{filter.Start, filter.Stop},
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			r.log.Error(err.Error())
		}
	}()

	var books entities.Books

	var book entities.Book
	var authorsUUID string
	var authors string

	for rows.Next() {
		err := rows.Scan(
			&book.ID,
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&authors,
			&authorsUUID,
		)
		if err != nil {
			return nil, err
		}
		var a entities.Authors
		aList := strings.Split(authors, ", ")
		aUUID := strings.Split(authorsUUID, ", ")
		for authorNumber, author := range aList {
			a = append(a, entities.Author{
				UUID:     aUUID[authorNumber],
				NickName: author,
			})
		}
		book.BookAuthors = a
		books.Books = append(books.Books, book)
	}
	return &books, err
}
