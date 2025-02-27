package bookusecase

import (
	"book/internal/dtos"
	"book/internal/entities"
	apperrors "book/internal/errors"
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
	filter, err := u.MapDTOToEntity(payload)
	if err != nil {
		return nil, err
	}

	booksCount, err := u.repo.BooksCount(ctx)
	if err != nil {
		return nil, err
	}

	if filter.Start > booksCount {
		return nil, apperrors.ErrPageNotFound
	}

	if filter.Stop > booksCount {
		filter.Stop = booksCount
	}

	filter.BooksCount = filter.Stop - filter.Start

	result, err := u.repo.Books(ctx, filter)
	if err != nil {
		return nil, err
	}

	return u.MapEntityToDTO(result)
}

func (u *getBooksUseCase) MapEntityToDTO(entity *entities.Books) (*dtos.BooksResponse, error) {
	response := &dtos.BooksResponse{}
	for _, book := range *entity {
		// get authors from entities.books and convert to dto
		var authors dtos.BookAuthorsResponse
		for _, author := range book.BookAuthors {
			authors.Authors = append(authors.Authors, dtos.BookAuthorResponse{
				UUID:     author.UUID.String,
				NickName: author.NickName.String,
			})
		}
		response.Books = append(response.Books, dtos.BookResponse{
			UUID:     book.UUID,
			Title:    book.Title,
			FilePath: book.BooksFileUUID.String,
			Authors:  authors,
		})
	}

	return response, nil
}

func (u *getBooksUseCase) MapDTOToEntity(dto *dtos.BooksRequest) (*entities.BooksFilter, error) {
	return &entities.BooksFilter{
		Start: (dto.Page - 1) * dto.Limit,
		Stop:  (dto.Page-1)*dto.Limit + dto.Limit,
	}, nil
}
