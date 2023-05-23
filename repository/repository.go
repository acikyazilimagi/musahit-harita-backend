package repository

import (
	"context"
	_ "embed"
	"github.com/Masterminds/squirrel"
	redisStore "github.com/acikkaynak/musahit-harita-backend/cache"
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/model/city"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/acikkaynak/musahit-harita-backend/utils/langutil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

var (
	psql                             = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	volunteerDistrictCountsTableName = "volunteer_district_counts"

	//go:embed tr-cities.json
	trCities []byte
	Cities   []city.City

	//go:embed tr-city-districts.json
	trDistricts []byte
	Districts   []city.District

	json = jsoniter.ConfigCompatibleWithStandardLibrary
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

	err = initDistricts()
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

	rawSql := psql.Select("district_id, count").From(volunteerDistrictCountsTableName)
	sql, args, err := rawSql.ToSql()
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var response feeds.Response
	feedResults := make([]feeds.Feed, 0)
	for rows.Next() {
		var feed feeds.Feed
		err := rows.Scan(&feed.DistrictId, &feed.VolunteerData)
		if err != nil {
			return nil, err
		}

		feedResults = append(feedResults, feed)
	}

	response.Count = len(feedResults)
	response.Results = feedResults

	return &response, nil
}

func initDistricts() error {
	err := json.Unmarshal(trCities, &Cities)
	if err != nil {
		return err
	}

	rd, err := redisStore.NewRedisStore()
	if err != nil {
		return err
	}

	var cityIdToCityName = make(map[int64]string)
	for _, c := range Cities {
		cityIdToCityName[c.Id] = c.Name
	}

	err = json.Unmarshal(trDistricts, &Districts)
	if err != nil {
		return err
	}

	for _, d := range Districts {
		cityName := langutil.ConvertTurkishCharsToEnglish(cityIdToCityName[d.CityId])
		districtName := langutil.ConvertTurkishCharsToEnglish(d.Name)
		key := strings.ToUpper(cityName + ":" + districtName)
		rd.SetKey(key, d.Id, 0)
	}

	return nil
}
