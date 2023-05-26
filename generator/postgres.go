package main

import (
	"context"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
)

func NewDB(conf Config) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(conf.DatabaseConnString)
	if err != nil {
		log.Logger().Fatal("Unable to parse DATABASE_URL", zap.Error(err))
		os.Exit(1)
	}

	cfg.MinConns = 5
	cfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Logger().Fatal("Unable to create connection pool", zap.Error(err))
		os.Exit(1)
	}

	return pool
}
