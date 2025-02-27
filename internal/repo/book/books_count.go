package bookrepo

import (
	"book/internal/adapters/storage"
	"context"
)

func (r *bookRepository) BooksCount(ctx context.Context) (int, error) {
	query := storage.Query{
		QueryName: "get books count",
		Query:     "select count(*) from books",
	}

	var booksCount int
	err := r.db.QueryRowContext(ctx, query).Scan(&booksCount)
	if err != nil {
		return 0, err
	}

	return booksCount, nil
}
