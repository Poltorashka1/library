package config

import "os"

type FileSystemConfig interface {
	StoragePath() string
}

type fileSystemConfig struct {
	storagePath string
}

func NewFileSystemConfig() FileSystemConfig {
	return &fileSystemConfig{storagePath: os.Getenv("FILE_STORAGE_PATH")}
}

func (fs *fileSystemConfig) StoragePath() string {
	return fs.storagePath
}
