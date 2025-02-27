package config

import (
	"os"
)

type DBConfig interface {
	DSN() string
}

type dbConfig struct {
	path string
}

func NewDbConfig() DBConfig {
	return &dbConfig{path: os.Getenv("DB_PATH")}
}

func (c *dbConfig) DSN() string {
	return c.path
}
