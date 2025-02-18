package bookusecase

import (
	"book/internal/dtos"
	"book/internal/logger"
	"book/internal/repo"
	"context"
)

type CreateBookUseCase interface {
	Run(ctx context.Context, payload *dtos.CreateBookRequest) (*dtos.CreateBookResponse, error)
}

type createBookUseCase struct {
	log  logger.Logger
	repo repo.Repository
}

func NewCreateBookUseCase(log logger.Logger, repo repo.Repository) CreateBookUseCase {
	return &createBookUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *createBookUseCase) Run(ctx context.Context, payload *dtos.CreateBookRequest) (*dtos.CreateBookResponse, error) {
	return nil, nil
}
