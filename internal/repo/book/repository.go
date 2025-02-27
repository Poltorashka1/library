package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	"book/internal/logger"
	"context"
)

type BookRepository interface {
	BookV3(ctx context.Context, uuid string) (*entities.Book, error)
	Books(ctx context.Context, payload *entities.BooksFilter) (*entities.Books, error)

	BooksCount(ctx context.Context) (int, error)
	BookFile(ctx context.Context, payload *entities.BookFileFilter) (*entities.BookFile, error)
	//BookAuthors(ctx context.Context, id int) (*entities.Authors, error)
}

type bookRepository struct {
	log logger.Logger
	db  storage.DB
	fs  storage.FS
}

func NewBookRepository(log logger.Logger, db storage.DB, fs storage.FS) BookRepository {
	return &bookRepository{log: log, db: db, fs: fs}
}
