package bookrepo

import (
	"book/internal/entities"
	"context"
)

// todo optimize []byte
func (r *bookRepository) BookFile(ctx context.Context, payload *entities.BookFileFilter) (*entities.BookFile, error) {
	file, err := r.fs.GET(ctx, payload.FileName, payload.FileType, payload.Chapter)
	if err != nil {
		return nil, err
	}

	return &entities.BookFile{File: file}, nil
}
