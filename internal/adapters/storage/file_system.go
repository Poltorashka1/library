package storage

import (
	"book/internal/config"
	apperrors "book/internal/errors"
	"book/internal/logger"
	"context"
	"errors"
	"io"
	"os"
)

type FS interface {
	GET(ctx context.Context, fileName string, fileType string, chapter string) ([]byte, error)
}

type fs struct {
	cfg    config.FileSystemConfig
	logger logger.Logger
}

func NewFileSystem(ctx context.Context, logger logger.Logger, cfg config.FileSystemConfig) FS {
	return &fs{
		cfg:    cfg,
		logger: logger,
	}
}

// todo refactor and optimize

func (fs *fs) GET(ctx context.Context, fileName string, fileType string, chapter string) ([]byte, error) {
	f := []byte(fileName)
	fileDir := fs.cfg.StoragePath() + "/" + fileType + "/" + string(f[2]) + "/" + string(f[3]) + "/" + fileName
	_, err := os.Stat(fileDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			/// todo feature return what problem does
			return nil, apperrors.ErrBookNotExist
		}
		return nil, err
	}
	var filePath string
	if chapter == "" {
		filePath = fileDir + "/" + fileName + "." + fileType
	} else {
		filePath = fileDir + "/" + chapter + "." + fileType
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return result, nil
}
