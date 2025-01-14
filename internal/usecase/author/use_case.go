package authorusecase

import (
	"book/internal/dtos"
	"book/internal/logger"
	"book/internal/repo"
)

type AuthorUseCase interface {
	Author(id int) *dtos.BookAuthorResponse
}

type authorUseCase struct {
	repo   repo.Repository
	logger logger.Logger
}

func NewAuthorUseCase(logger logger.Logger, repo repo.Repository) AuthorUseCase {
	return &authorUseCase{repo: repo, logger: logger}
}
