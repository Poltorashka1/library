package config

import (
	"fmt"
	"os"
)

type HttpConfig interface {
	Addr() string
}

type httpServerConfig struct {
	host string
	port string
}

func LoadHttpConfig() HttpConfig {
	return &httpServerConfig{host: os.Getenv("HTTP_HOST"), port: os.Getenv("HTTP_PORT")}
}

func (c *httpServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
