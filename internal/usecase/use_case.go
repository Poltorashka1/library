package usecase

import (
	"book/internal/logger"
	"book/internal/repo"
	"book/internal/usecase/author"
	"book/internal/usecase/book"
)

//go:generate mockery --name=UseCase

type UseCase interface {
	bookusecase.BookUseCase
	authorusecase.AuthorUseCase
}

type useCase struct {
	bookusecase.BookUseCase
	authorusecase.AuthorUseCase
}

func NewUseCase(logger logger.Logger, repo repo.Repository) UseCase {
	return &useCase{
		bookusecase.NewBookUseCase(logger, repo),
		authorusecase.NewAuthorUseCase(logger, repo),
	}
}
