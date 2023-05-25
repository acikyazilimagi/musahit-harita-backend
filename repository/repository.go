package repository

import (
	"context"
	_ "embed"
	"github.com/acikkaynak/musahit-harita-backend/utils/stringutils"
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
	trCities                                     []byte
	CityIdToMap                                  = make(map[int]model.City)
	DistrictIdToMap                              = make(map[int]model.District)
	NeighborhoodIdToMap                          = make(map[int]model.Neighborhood)
	CityToDistrictToNeighborhoodToNeighborhoodId = make(map[string]map[string]map[string]int)
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
		WHERE neighbourhood_id = $1
		LIMIT 1
	  ), box_numbers AS (
		SELECT b.id AS building_id, ARRAY_AGG(bb.box_no ORDER BY bb.box_no) AS box_numbers
		FROM volunteer_counts vc
		JOIN selected_neighbourhood sn ON vc.location_id = sn.id
		LEFT JOIN buildings b ON vc.building_id = b.id
		LEFT JOIN ballot_boxes bb ON b.id = bb.building_id
		GROUP BY b.id
	  )
	  SELECT b.name AS building_name, bn.box_numbers
	  FROM buildings b
	  LEFT JOIN box_numbers bn ON b.id = bn.building_id;	  
	`
	args := []interface{}{neighborhoodId}
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var response feeds.FeedDetailResponse
	feedDetailResults := make([]feeds.FeedDetail, 0)
	for rows.Next() {
		var feedDetail feeds.FeedDetail
		err := rows.Scan(&feedDetail.BuildingName, &feedDetail.BallotBoxNos)
		if err != nil {
			return nil, err
		}
		feedDetailResults = append(feedDetailResults, feedDetail)
	}
	response.Details = feedDetailResults
	response.NeighborhoodId = neighborhoodId
	response.Intensity = rand.Intn(5) + 1                                                                      // change to real intensity
	response.LastUpdateTime = time.Now().Add(-time.Minute * time.Duration(rand.Intn(60))).Format(time.RFC3339) // change to real time
	return &response, nil
}

func (r *Repository) GetFeedDetailFromMemory(neighborhoodId int) (*feeds.FeedDetailResponse, error) {
	var response feeds.FeedDetailResponse

	response.NeighborhoodId = neighborhoodId
	response.LastUpdateTime = OvoBuildingStore.LastUpdateTime.Format(time.RFC3339)

	response.Intensity = OvoBuildingStore.NeighborhoodIdToAvgScore[neighborhoodId]
	if response.Intensity == 0 {
		response.Intensity = 1
	}

	feedDetailResults := make([]feeds.FeedDetail, 0)
	for _, building := range OvoBuildingStore.NeighborhoodIdToBuildings[neighborhoodId] {
		var feedDetail feeds.FeedDetail
		feedDetail.BuildingName = building.BuildingName
		// TODO: change to real ballot box numbers
		feedDetail.BallotBoxNos = []int{}
		feedDetailResults = append(feedDetailResults, feedDetail)
	}
	response.Details = feedDetailResults
	return &response, nil
}

func (r *Repository) ApplyVolunteer(volunteer model.VolunteerDoc) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var volunteerId int
	rawSql := psql.Select("id").From("volunteers").Where("volunteer_doc->>'email' = ?", volunteer.Email).Limit(1)
	sql, args, err := rawSql.ToSql()
	if err != nil {
		return 0, err
	}
	row, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	for row.Next() {
		err = row.Scan(&volunteerId)
		if err != nil {
			return 0, err
		}
	}
	if volunteerId == 0 {
		var locationId int
		rawSql := psql.Select("id").From("locations").Where("neighbourhood_id = ?", volunteer.NeighborhoodId).Limit(1)
		sql, args, err := rawSql.ToSql()
		lrow, err := r.pool.Query(ctx, sql, args...)
		if err != nil {
			return 0, err
		}
		for lrow.Next() {
			err = lrow.Scan(&locationId)
			if err != nil {
				return 0, err
			}
		}
		volunteerDoc, err := json.Marshal(volunteer)
		if err != nil {
			return 0, err
		}
		volunteer := model.Volunteer{
			VolunteerDoc: volunteerDoc,
			LocationId:   locationId,
			Confirmed:    false,
			SourceId:     2,
			BuildingId:   0,
		}
		insertSql := psql.Insert("volunteers").Columns("volunteer_doc", "location_id", "confirmed", "source_id", "building_id").Values(volunteer.VolunteerDoc, volunteer.LocationId, volunteer.Confirmed, volunteer.SourceId, volunteer.BuildingId).Suffix("RETURNING id")
		sql, args, err = insertSql.ToSql()
		if err != nil {
			return 0, err
		}
		irow, err := r.pool.Query(ctx, sql, args...)
		if err != nil {
			return 0, err
		}
		for irow.Next() {
			err = irow.Scan(&volunteerId)
			if err != nil {
				return 0, err
			}
		}
		return volunteerId, nil
	}
	return volunteerId, nil
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

func (r *Repository) GetFeedsFromMemory() (*feeds.Response, error) {
	response := make([]feeds.Feed, 0)
	ovoBuildingStore := OvoBuildingStore

	for k, v := range ovoBuildingStore.NeighborhoodIdToAvgScore {
		response = append(response, feeds.Feed{
			NeighborhoodId: k,
			VolunteerData:  v,
		})
	}

	return &feeds.Response{
		Count:   len(response),
		Results: response,
	}, nil
}

func initCities() error {
	cityIdToMap := make(map[string]model.City)
	err := json.Unmarshal(trCities, &cityIdToMap)
	if err != nil {
		return err
	}

	for _, city := range cityIdToMap {
		CityIdToMap[city.Id] = city
		CityToDistrictToNeighborhoodToNeighborhoodId[city.Name] = make(map[string]map[string]int)
		for _, district := range city.Districts {
			DistrictIdToMap[district.Id] = district
			CityToDistrictToNeighborhoodToNeighborhoodId[city.Name][district.Name] = make(map[string]int)
			for _, neighborhood := range district.Neighborhoods {
				NeighborhoodIdToMap[neighborhood.Id] = neighborhood
				// It is possible that a neighborhood name has parentheses in it. We need to parse them.
				for _, parsedName := range stringutils.ParseParentheses(neighborhood.Name) {
					CityToDistrictToNeighborhoodToNeighborhoodId[city.Name][district.Name][parsedName] = neighborhood.Id
				}
			}
		}
	}

	return nil
}
