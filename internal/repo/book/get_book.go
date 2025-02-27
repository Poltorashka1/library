package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	apperrors "book/internal/errors"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// todo get image

const getBookByUUID = `select t1.uuid,
	  			 t1.isbn,
	  			 t1.title,
	  			 t1.publication_year,
	  			 t1.description,
	  			 t1.books_file_uuid,
	  			 t1.publisher,
	  			 t3.uuid,
	  			 t3.nickname
			from books as t1
	    		left join book_author as t2 on t2.book_id = t1.id
	    		left join authors as t3 on t2.author_id = t3.id
			WHERE t1.uuid = $1;`

// Book get book by uuid, err return apperrors.ErrBookNotFound
func (r *bookRepository) Book(ctx context.Context, uuid string) (*entities.Book, error) {
	query := storage.Query{
		QueryName: "get book by uuid",
		Query:     getBookByUUID,
		Args:      []any{uuid},
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

	book := &entities.Book{}

	if !rows.Next() {
		return nil, apperrors.ErrBookNotFound
	}

	for {
		var author entities.BookAuthor
		err := rows.Scan(
			&book.UUID,
			&book.ISBN,
			&book.Title,
			&book.PublicationYear,
			&book.Description,
			&book.BooksFileUUID,
			&book.Publisher,
			&author.UUID,
			&author.NickName,
		)
		if err != nil {
			return nil, err
		}
		book.BookAuthors = append(book.BookAuthors, author)

		if !rows.Next() {
			break
		}
	}

	return book, err
}

func (r *bookRepository) BookV2(ctx context.Context, uuid string) (*entities.Book, error) {
	query := storage.Query{
		QueryName: "get book by uuid",
		Query: `SELECT uuid,
   					 isbn,
   					 title,
  					 publication_year,
   					 description,
   	   				 books_file_uuid,
   					 publisher 
  			  FROM books 
  			WHERE uuid = $1;`,
		Args: []any{uuid},
	}

	row := r.db.QueryRowContext(ctx, query)

	book := &entities.Book{}

	err := row.Scan(&book.UUID, &book.ISBN, &book.Title, &book.PublicationYear, &book.Description, &book.BooksFileUUID, &book.Publisher)
	if err != nil {
		return nil, err
	}

	query2 := storage.Query{
		QueryName: "get book authors by book uuid",
		Query: `SELECT t3.uuid,
       					t3.nickname
FROM authors AS t3
JOIN book_author AS t2 ON t2.author_id = t3.id
JOIN books AS t1 ON t2.book_id = t1.id
WHERE t1.uuid = $1;`,
		Args: []any{uuid},
	}

	rows, err := r.db.QueryContext(ctx, query2)
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
		return nil, apperrors.ErrBookNotFound
	}

	for {
		var author entities.BookAuthor
		err = rows.Scan(&author.UUID, &author.NickName)
		if err != nil {
			return nil, err
		}

		book.BookAuthors = append(book.BookAuthors, author)

		if !rows.Next() {
			break
		}
	}

	return book, nil
}

func (r *bookRepository) BookV3(ctx context.Context, uuid string) (*entities.Book, error) {
	query := storage.Query{
		QueryName: "get book by uuid",
		Query: `select t1.uuid,
	  			 t1.isbn,
	  			 t1.title,
	  			 t1.publication_year,
	  			 t1.description,
	  			 t1.books_file_uuid,
	  			 t1.publisher,
	  			 GROUP_CONCAT(t3.nickname, ', ') AS authors,
       				GROUP_CONCAT(t3.uuid, ', ')     AS authors_uuid
			from books as t1
	    		left join book_author as t2 on t2.book_id = t1.id
	    		left join authors as t3 on t2.author_id = t3.id
			WHERE t1.uuid = $1
GROUP BY t1.id;`,
		Args: []any{uuid},
	}

	row := r.db.QueryRowContext(ctx, query)

	book := &entities.Book{}

	var nicknameStr string
	var uuidStr string
	err := row.Scan(
		&book.UUID,
		&book.ISBN,
		&book.Title,
		&book.PublicationYear,
		&book.Description,
		&book.BooksFileUUID,
		&book.Publisher,
		&nicknameStr,
		&uuidStr,
	)
	if err != nil {
		return nil, err
	}

	nL := strings.Split(nicknameStr, ", ")
	uL := strings.Split(uuidStr, ", ")
	if len(nL) != len(uL) {
		return nil, fmt.Errorf("mismatch between author nicknames and UUIDs")
	}

	book.BookAuthors = make([]entities.BookAuthor, 0, len(nL))
	for i, v := range nL {
		var author entities.BookAuthor
		author.NickName = sql.NullString{String: v}
		author.UUID = sql.NullString{String: uL[i]}
		book.BookAuthors = append(book.BookAuthors, author)
	}

	return book, err
}
