package main

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config stores app configuration
type Config struct {
	DBPath            string `env:"DB_PATH,default=my.bdb"`
	HTTPServerAddress string `env:"HTTP_SERVER_ADDRESS,default=0.0.0.0:8080"`
	ReadTimeout       int    `env:"READ_TIMEOUT,default=5"`
	IdleTimeout       int    `env:"IDLE_TIMEOUT,default=30"`
	ShutdownTimeout   int    `env:"SHUTDOWN_TIMEOUT,default=10"`
}

// NewConfig reads config from env and creates config struct
func NewConfig() (*Config, error) {
	ctx := context.Background()
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
