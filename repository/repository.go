package repository

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/model"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	psql                             = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	volunteerLocationCountsTableName = "volunteer_counts"

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

func (r *Repository) GetFeedDetail(neighborhoodId int) (*feeds.FeedDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	sql := `
	WITH selected_neighbourhood AS (
		SELECT id
		FROM locations
		WHERE neighbourhood_ysk_id = $1
		LIMIT 1
	  )
	  SELECT b.name AS building_name, STRING_AGG(CAST(bb.box_no AS VARCHAR), ' - ') AS combined_box_no
	  FROM volunteer_counts vc
	  JOIN selected_neighbourhood sn ON vc.location_id = sn.id
	  LEFT JOIN buildings b ON vc.building_id = b.id
	  LEFT JOIN ballot_boxes bb ON b.id = bb.building_id
	  GROUP BY b.id;
	`
	args := []interface{}{neighborhoodId}
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var response feeds.FeedDetailResponse
	feedDetailResults := make([]string, 0)
	for rows.Next() {
		var feedDetail feeds.FeedDetail
		err := rows.Scan(&feedDetail.BuildingName, &feedDetail.BallotBoxCombine)
		if err != nil {
			return nil, err
		}
		str := fmt.Sprintf("%s - %s", feedDetail.BuildingName, feedDetail.BallotBoxCombine)
		feedDetailResults = append(feedDetailResults, str)
	}
	response.Details = feedDetailResults
	response.NeighborhoodId = neighborhoodId
	response.Intensity = rand.Intn(5) + 1                                                                      // change to real intensity
	response.LastUpdateTime = time.Now().Add(-time.Minute * time.Duration(rand.Intn(60))).Format(time.RFC3339) // change to real time
	return &response, nil
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
