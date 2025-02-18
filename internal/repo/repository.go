package repo

import (
	"book/internal/adapters/storage"
	"book/internal/logger"
	"book/internal/repo/author"
	"book/internal/repo/book"
)

type Repository interface {
	bookrepo.BookRepository
	authorrepo.AuthorRepository
}

// todo mb storage.db сюда
type repository struct {
	bookrepo.BookRepository
	authorrepo.AuthorRepository
}

func NewRepository(log logger.Logger, db storage.DB, fs storage.FS) Repository {
	return &repository{
		bookrepo.NewBookRepository(log, db, fs),
		authorrepo.NewAuthorRepository(log, db),
	}
}
