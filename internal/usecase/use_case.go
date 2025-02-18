package usecase

import (
	"book/internal/logger"
	"book/internal/repo"
	authorusecase "book/internal/usecase/author"
	"book/internal/usecase/book"
)

//go:generate mockery --name=UseCase

type UseCase struct {
	*BookUseCase
	*AuthorUseCase
}

func NewUseCase(log logger.Logger, repo repo.Repository) *UseCase {
	return &UseCase{NewBookUseCase(log, repo),
		NewAuthorUseCase(log, repo)}
}

type BookUseCase struct {
	bookusecase.CreateBookUseCase
	bookusecase.GetBookUseCase
	bookusecase.GetBooksUseCase
}

func NewBookUseCase(log logger.Logger, repo repo.Repository) *BookUseCase {
	return &BookUseCase{
		bookusecase.NewCreateBookUseCase(log, repo),
		bookusecase.NewGetBookUseCase(log, repo),
		bookusecase.NewGetBooksUseCase(log, repo),
	}
}

type AuthorUseCase struct {
	authorusecase.GetAuthorUseCase
}

func NewAuthorUseCase(log logger.Logger, repo repo.Repository) *AuthorUseCase {
	return &AuthorUseCase{
		authorusecase.NewGetAuthorUseCase(log, repo),
	}
}
