package bookusecase

import (
	"book/internal/dtos"
	"book/internal/entities"
	"book/internal/logger"
	"book/internal/repo"
	"context"
)

type GetBooksUseCase interface {
	Run(ctx context.Context, payload *dtos.BooksRequest) (*dtos.BooksResponse, error)
}

// todo add error if book file not found
type getBooksUseCase struct {
	log  logger.Logger
	repo repo.Repository
}

func NewGetBooksUseCase(log logger.Logger, repo repo.Repository) GetBooksUseCase {
	return &getBooksUseCase{log, repo}
}
func (u *getBooksUseCase) Run(ctx context.Context, payload *dtos.BooksRequest) (*dtos.BooksResponse, error) {
	filter := entities.BookFilter{
		Start: (payload.Page - 1) * payload.Limit,
		Stop:  (payload.Page-1)*payload.Limit + payload.Limit,
	}

	// get books list with authors
	result, err := u.repo.Books(ctx, filter)
	if err != nil {
		return nil, err
	}

	// todo add converter
	// convert books to dto
	var response dtos.BooksResponse
	for _, book := range result.Books {
		// get authors from entities.books and convert to dto
		var authors dtos.BookAuthorsResponse
		for _, author := range book.BookAuthors {
			authors.Authors = append(authors.Authors, dtos.BookAuthorResponse{
				UUID:     author.UUID,
				NickName: author.NickName,
			})
		}

		// todo сделать что то с file path
		// convert book from entities.Books and convert to dto
		//if book.BooksFileUUID == nil {
		//	*book.BooksFileUUID = ""
		//}
		response.Books = append(response.Books, dtos.BookResponse{
			UUID:     book.UUID,
			Title:    book.Title,
			FilePath: book.BooksFileUUID.String,
			Authors:  authors,
		})
	}

	return &response, nil
}
