package bookusecase

import (
	"book/internal/dtos"
	"book/internal/entities"
	"book/internal/logger"
	"book/internal/repo"
	"context"
)

type GetBookUseCase interface {
	Run(ctx context.Context, uuid string) (*dtos.BookResponse, error)
}

// todo add error if book file not found
type getBookUseCase struct {
	log  logger.Logger
	repo repo.Repository
}

func NewGetBookUseCase(log logger.Logger, repo repo.Repository) GetBookUseCase {
	return &getBookUseCase{log, repo}
}

func (u *getBookUseCase) Run(ctx context.Context, uuid string) (*dtos.BookResponse, error) {
	result, err := u.repo.BookV3(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return u.MapEntityToDTO(result)
}

func (u *getBookUseCase) MapEntityToDTO(entity *entities.Book) (*dtos.BookResponse, error) {
	var authors dtos.BookAuthorsResponse
	for _, author := range entity.BookAuthors {
		authors.Authors = append(authors.Authors, dtos.BookAuthorResponse{
			UUID:     author.UUID.String,
			NickName: author.NickName.String,
		})
	}

	return &dtos.BookResponse{
		ISBN:            entity.ISBN,
		UUID:            entity.UUID,
		Title:           entity.Title,
		PublicationYear: entity.PublicationYear,
		Description:     entity.Description,
		Authors:         authors,
	}, nil
}
