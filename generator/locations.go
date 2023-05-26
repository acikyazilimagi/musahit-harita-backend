package main

import (
	"context"
	_ "embed"
	"github.com/Masterminds/squirrel"
	"github.com/acikkaynak/musahit-harita-backend/model"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type Location struct {
	CityID           int64  `json:"city_id"`
	CityName         string `json:"city_name"`
	DistrictID       int64  `json:"district_id"`
	DistrictName     string `json:"district_name"`
	NeighborhoodID   int64  `json:"neighborhood_id"`
	NeighborhoodName string `json:"neighborhood_name"`
	CityCode         string `json:"city_code"`
	Lat              string `json:"lat"`
	Long             string `json:"long"`
}

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func Migrate(pool *pgxpool.Pool) {
	rawSql := psql.Select("city_id",
		"city_name",
		"district_id",
		"district_name",
		"neighborhood_id",
		"neighborhood_name",
		"latitude",
		"longitude").From("locations").OrderBy("id ASC")
	sql, _, err := rawSql.ToSql()
	if err != nil {
		log.Logger().Fatal("Unable to generate sql", zap.Error(err))
	}

	rows, err := pool.Query(context.Background(), sql)
	if err != nil {
		log.Logger().Fatal("Unable to query", zap.Error(err))
	}

	response := map[int64]*model.City{}
	lastCity := &model.City{}
	lastDistrict := &model.District{}
	cityDistricts := make(map[int64][]model.District)
	cityDistrictNeighborhoods := make(map[int64]map[int64][]model.Neighborhood)
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.CityID,
			&location.CityName,
			&location.DistrictID,
			&location.DistrictName,
			&location.NeighborhoodID,
			&location.NeighborhoodName,
			&location.Lat,
			&location.Long)
		if err != nil {
			log.Logger().Fatal("Unable to scan", zap.Error(err))
		}

		city, ok := response[location.CityID]
		if !ok {
			city = &model.City{
				Id:        int(location.CityID),
				Name:      location.CityName,
				Districts: make([]model.District, 0),
				Type:      "city",
			}
			response[location.CityID] = city
			cityDistricts[location.CityID] = make([]model.District, 0)
			cityDistrictNeighborhoods[location.CityID] = make(map[int64][]model.Neighborhood)
		}

		if lastCity == nil || lastCity.Id == 0 {
			lastCity = city
			lastDistrict = &model.District{
				Id:   int(location.DistrictID),
				Name: location.DistrictName,
				Type: "district",
			}
		}

		if lastCity.Id != int(location.CityID) {
			response[int64(lastCity.Id)].Districts = cityDistricts[int64(lastCity.Id)]
			//fmt.Println("cityDistricts[lastCityId]", cityDistricts[int64(lastCity.Id)])
			cityDistricts[location.CityID] = make([]model.District, 0)
			cityDistrictNeighborhoods[location.CityID] = make(map[int64][]model.Neighborhood)
		}

		if lastDistrict.Id != int(location.DistrictID) {
			//neighborhoods := cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)]
			//fmt.Println("cityDistrictNeighborhoods[lastCityId][lastDistrictId]", cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)])

			lastDistrict.Neighborhoods = cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)]
			cityDistricts[int64(lastCity.Id)] = append(cityDistricts[int64(lastCity.Id)], *lastDistrict)
			response[int64(lastCity.Id)].Districts = cityDistricts[int64(lastCity.Id)]
		}

		Lat, err := strconv.ParseFloat(location.Lat, 64)
		if err != nil {
			log.Logger().Fatal("Unable to parse lat", zap.Error(err))
		}
		Long, err := strconv.ParseFloat(location.Long, 64)
		if err != nil {
			log.Logger().Fatal("Unable to parse long", zap.Error(err))
		}
		neighborhood := &model.Neighborhood{
			Id:         int(location.NeighborhoodID),
			Name:       location.NeighborhoodName,
			DistrictId: int(location.DistrictID),
			Type:       "neighborhood",
			Geo: &model.Geo{
				Lat:  Lat,
				Long: Long,
			},
		}
		cityDistrictNeighborhoods[location.CityID][location.DistrictID] = append(cityDistrictNeighborhoods[location.CityID][location.DistrictID], *neighborhood)

		lastCity = city
		lastDistrict = &model.District{
			Id:     int(location.DistrictID),
			Name:   location.DistrictName,
			CityID: int(location.CityID),
			Type:   "district",
		}
	}

	lastDistrict.Neighborhoods = cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)]
	cityDistricts[int64(lastCity.Id)] = append(cityDistricts[int64(lastCity.Id)], *lastDistrict)
	response[int64(lastCity.Id)].Districts = cityDistricts[int64(lastCity.Id)]

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Logger().Fatal("Unable to marshal", zap.Error(err))
	}
	err = os.WriteFile("tr_election_locations_2023.json", bytes, 0644)
	if err != nil {
		log.Logger().Fatal("Unable to write file", zap.Error(err))
	}
}
