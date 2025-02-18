package authorusecase

import (
	"book/internal/logger"
	"book/internal/repo"
	"context"
)

type GetAuthorUseCase interface {
	Run(ctx context.Context, id int) error
}

type getAuthorUseCase struct {
	log  logger.Logger
	repo repo.Repository
}

func NewGetAuthorUseCase(log logger.Logger, repo repo.Repository) GetAuthorUseCase {
	return &getAuthorUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *getAuthorUseCase) Run(ctx context.Context, id int) error {
	return nil
}
