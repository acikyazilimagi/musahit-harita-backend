package main

import (
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	DatabaseConnString string `env:"DATABASE_CONN_STRING" envDefault:"postgres://postgres:postgres@localhost:5435/postgres?sslmode=disable"`
}

func ParseConfig() Config {
	var conf Config
	err := env.Parse(&conf)
	if err != nil {
		log.Logger().Fatal("Unable to parse config", zap.Error(err))
		os.Exit(1)
	}
	return conf
}
