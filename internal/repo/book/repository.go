package bookrepo

import (
	"book/internal/adapters/storage"
	"book/internal/entities"
	"book/internal/logger"
	"context"
)

type BookRepository interface {
	Book(ctx context.Context, uuid string) (*entities.Book, error)
	Books(ctx context.Context, payload entities.BookFilter) (*entities.Books, error)
	BookFile(ctx context.Context, payload entities.BookFileFilter) (*entities.BookFile, error)
	//BookAuthors(ctx context.Context, id int) (*entities.Authors, error)
}

type bookRepository struct {
	db     storage.DB
	fs     storage.FS
	logger logger.Logger
}

func NewBookRepository(logger logger.Logger, db storage.DB, fs storage.FS) BookRepository {
	return &bookRepository{db: db, logger: logger, fs: fs}
}
