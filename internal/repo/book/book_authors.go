package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	"context"
)

// todo id to uuid

func (r *bookRepository) BookAuthors(ctx context.Context, id int) (*entities.Authors, error) {
	query := storage.Query{
		QueryName: "get book authors",
		Query: `SELECT a.nickname,a.name, a.surname, a.patronymic
			FROM authors a
			JOIN book_author ba ON a.id = ba.author_id
			JOIN books b ON ba.book_id = b.id
			WHERE b.id = $1`,
		Args: []any{id},
	}

	authors, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var bookAuthors entities.Authors
	for authors.Next() {
		var author entities.Author
		err = authors.Scan(
			&author.Name,
			&author.Surname,
			&author.Patronymic,
		)
		if err != nil {
			return nil, err
		}
		bookAuthors.Authors = append(bookAuthors.Authors, author)
	}
	return &bookAuthors, nil

}
