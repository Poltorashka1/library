package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	apperrors "book/internal/errors"
	"context"
)

// todo get image

// Book get book by uuid, err return apperrors.ErrBookNotFound
func (r *bookRepository) Book(ctx context.Context, uuid string) (*entities.Book, error) {
	query := storage.Query{
		QueryName: "get book by uuid",
		Query: `select t1.uuid,
	  			 t1.isbn,
	  			 t1.title,
	  			 t1.publication_year,
	  			 t1.description,
	  			 t1.books_file_uuid,
	  			 t1.publisher,
	  			 t3.uuid,
	  			 t3.nickname,
	  			 t3.name,
	  			 t3.surname,
	  			 t3.patronymic
			from books as t1
	    		join book_author as t2 on t2.book_id = t1.id
	    		join authors as t3 on t2.author_id = t3.id
			WHERE t1.uuid = $1;`,
		Args: []any{uuid},
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

	// todo optimize при нескольких авторах происходит перезаписывание информации о книге
	book := new(entities.Book)

	if !rows.Next() {
		return nil, apperrors.ErrBookNotFound
	}

	for {
		var author entities.Author
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
			&author.Name,
			&author.Surname,
			&author.Patronymic,
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
