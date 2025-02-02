package config

import "os"

type FileSystemConfig interface {
	StoragePath() string
}

type fileSystemConfig struct {
	storagePath string
}

func NewFileSystemConfig() FileSystemConfig {
	return &fileSystemConfig{storagePath: os.Getenv("BOOKS_FILE_DIR")}
}

func (fs *fileSystemConfig) StoragePath() string {
	return fs.storagePath
}
