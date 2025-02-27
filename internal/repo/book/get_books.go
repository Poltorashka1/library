package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	apperrors "book/internal/errors"
	"context"
	"database/sql"
	"strings"
)

// Books get books with authors, genres; err return apperrors.ErrPageNotFound
func (r *bookRepository) Books(ctx context.Context, filter *entities.BooksFilter) (*entities.Books, error) {
	query := storage.Query{
		QueryName: "select books with pagination",
		Query: `SELECT t1.id,t1.uuid,t1.isbn,t1.title,t1.publication_year,t1.description,
       				GROUP_CONCAT(t3.nickname, ', ') AS authors,
       				GROUP_CONCAT(t3.uuid, ', ')     AS authors_uuid
				FROM books AS t1
         			Left JOIN book_author AS t2 ON t2.book_id = t1.id
         			LEft JOIN authors AS t3 ON t2.author_id = t3.id
				WHERE t1.id > $1 and t1.id <= $2
				GROUP BY t1.id
				ORDER BY t1.id;`,
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

	if !rows.Next() {
		return nil, apperrors.ErrPageNotFound
	}

	books := make(entities.Books, 0, filter.BooksCount)

	for {
		var book entities.Book
		var authorsUUID sql.NullString
		var authorsNickName sql.NullString
		err := rows.Scan(
			&book.ID,
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&authorsNickName,
			&authorsUUID,
		)
		if err != nil {
			return nil, err
		}

		var authors entities.BookAuthors

		if !authorsUUID.Valid || !authorsNickName.Valid {
			authors = append(authors, entities.BookAuthor{
				NickName: sql.NullString{String: "No information"},
			})
		} else {
			aList := strings.Split(authorsNickName.String, ", ")
			aUUID := strings.Split(authorsUUID.String, ", ")

			// todo if len(aList) != len(aUUID)
			for authorNumber, nickname := range aList {
				authors = append(authors, entities.BookAuthor{
					UUID:     sql.NullString{String: aUUID[authorNumber]},
					NickName: sql.NullString{String: nickname},
				})
			}
		}

		book.BookAuthors = authors
		books = append(books, book)

		if !rows.Next() {
			break
		}
	}

	return &books, err
}

const getBooksWithPagination = `
				SELECT t1.id,t1.uuid,t1.isbn,t1.title,t1.publication_year,t1.description,
       				t3.nickname,
       				t3.uuid
				FROM books AS t1
         			LEFT JOIN book_author AS t2 ON t2.book_id = t1.id
         			LEFT JOIN authors AS t3 ON t2.author_id = t3.id
				WHERE t1.id > $1 and t1.id <= $2
				ORDER BY t1.id`

// BooksV2 get books with authors; err return apperrors.ErrPageNotFound
func (r *bookRepository) BooksV2(ctx context.Context, filter *entities.BooksFilter) (*entities.Books, error) {
	query := storage.Query{
		QueryName: "select books with pagination",
		Query:     getBooksWithPagination,
		Args:      []any{filter.Start, filter.Stop},
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

	if !rows.Next() {
		return nil, apperrors.ErrPageNotFound
	}

	books := make(entities.Books, 0, filter.BooksCount)

	bookMap := map[int]entities.Book{}
	for {
		var book entities.Book

		var author entities.BookAuthor
		err := rows.Scan(
			&book.ID,
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&author.NickName,
			&author.UUID,
		)
		if err != nil {
			return nil, err
		}
		// todo что если два невалидный автора, тогда "No information" запишется дважды
		if !author.UUID.Valid || !author.NickName.Valid {
			author.NickName.String = "No information"
		}

		if existBook, exist := bookMap[book.ID]; !exist {
			book.BookAuthors = append(book.BookAuthors, author)
			bookMap[book.ID] = book
		} else {
			existBook.BookAuthors = append(existBook.BookAuthors, author)
			bookMap[book.ID] = existBook
		}

		if !rows.Next() {
			break
		}
	}

	// todo sorting disappears
	for _, book := range bookMap {
		books = append(books, book)
	}

	return &books, err
}
