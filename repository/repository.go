package repository

import (
	"context"
	_ "embed"
	"github.com/Masterminds/squirrel"
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/model"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	psql                             = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	volunteerLocationCountsTableName = "volunteer_locations_counts"

	//go:embed city-district-neighborhood.json
	trCities            []byte
	CityIdToMap         = make(map[int]model.City)
	DistrictIdToMap     = make(map[int]model.District)
	NeighborhoodIdToMap = make(map[int]model.Neighborhood)
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Close()
}

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

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Logger().Fatal("Unable to create connection pool", zap.Error(err))
		os.Exit(1)
	}

	err = initCities()
	if err != nil {
		log.Logger().Fatal("Unable to initialize districts", zap.Error(err))
		os.Exit(1)
	}

	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Close() {
	r.pool.Close()
}

func (r *Repository) GetFeeds() (*feeds.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rawSql := psql.Select("neighbourhood_id, count").From(volunteerLocationCountsTableName)
	sql, args, err := rawSql.ToSql()
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var response feeds.Response
	feedResults := make([]feeds.Feed, 0)
	for rows.Next() {
		var feed feeds.Feed
		err := rows.Scan(&feed.NeighborhoodId, &feed.VolunteerData)
		if err != nil {
			return nil, err
		}

		feedResults = append(feedResults, feed)
	}

	response.Count = len(feedResults)
	response.Results = feedResults

	return &response, nil
}

func initCities() error {
	cityIdToMap := make(map[string]model.City)
	err := json.Unmarshal(trCities, &cityIdToMap)
	if err != nil {
		return err
	}

	for _, city := range cityIdToMap {
		CityIdToMap[city.Id] = city
		for _, district := range city.Districts {
			DistrictIdToMap[district.Id] = district
			for _, neighborhood := range district.Neighborhoods {
				NeighborhoodIdToMap[neighborhood.Id] = neighborhood
			}
		}
	}

	return nil
}
