package repository

import (
	"context"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"os"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New() *Repository {
	dbUrl := os.Getenv("DB_CONNECTION_STRING")
	cfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Logger().Fatal("Unable to parse DATABASE_URL", zap.Error(err))
		os.Exit(1)
	}

	cfg.MinConns = 5
	cfg.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Logger().Fatal("Unable to create connection pool", zap.Error(err))
		os.Exit(1)
	}

	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Close() {
	r.pool.Close()
}
