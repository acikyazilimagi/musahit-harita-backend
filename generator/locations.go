package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/acikkaynak/musahit-harita-backend/model"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
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

func UpdateGeolocation(pool *pgxpool.Pool) {
	// decode fix.json to map
	fixFile, err := os.ReadFile("generator/fix.json")
	if err != nil {
		log.Logger().Fatal("Unable to read fix.json", zap.Error(err))
	}

	type Fix struct {
		CityID         int     `json:"cityID"`
		DistrictID     int     `json:"districtID"`
		NeighborhoodID int     `json:"neighborhoodID"`
		Neighborhood   string  `json:"neighborhood"`
		Lat            float64 `json:"lat"`
		Lng            float64 `json:"lng"`
	}
	type All struct {
		CityID           int    `json:"cityID"`
		DistrictID       int    `json:"districtID"`
		NeighborhoodID   int    `json:"neighborhoodID"`
		Neighborhood     string `json:"neighborhood"`
		CityNameOriginal string `json:"cityNameOriginal"`
		CityNameFetched  string `json:"cityNameFetched"`
		Lat              string `json:"lat"`
		Lng              string `json:"lng"`
	}

	var fix []Fix
	err = json.Unmarshal(fixFile, &fix)
	if err != nil {
		log.Logger().Fatal("Unable to unmarshal fix.json", zap.Error(err))
	}

	neighIds := make([]int, 0)
	//
	fixMap := make(map[int]Fix)
	fixCount := 0
	for _, f := range fix {
		fixCount++
		neighIds = append(neighIds, f.NeighborhoodID)
		//
		fixMap[f.NeighborhoodID] = f
		//latStr := strconv.FormatFloat(f.Lat, 'f', -1, 64)
		//lngStr := strconv.FormatFloat(f.Lng, 'f', -1, 64)
		//rawSql := psql.Update("neighborhood").Set("lat", latStr).Set("lng", lngStr).Where(squirrel.Eq{"id": f.NeighborhoodID})
		//sql, args, err := rawSql.ToSql()
		//if err != nil {
		//	log.Logger().Fatal("Unable to generate sql", zap.Error(err))
		//}
		//
		//fmt.Println(sql)
		//
		//_, err = pool.Exec(context.Background(), sql, args...)
		//if err != nil {
		//	log.Logger().Fatal("Unable to exec", zap.Error(err))
		//}
		//
		//log.Logger().Info("Updated", zap.Int("neighborhood_id", f.NeighborhoodID), zap.Float64("lat", f.Lat), zap.Float64("lng", f.Lng))
	}
	fmt.Println("fixCount", fixCount)

	allFile, err := os.ReadFile("generator/all.json")
	if err != nil {
		log.Logger().Fatal("Unable to read all.json", zap.Error(err))
	}

	var all []All
	err = json.Unmarshal(allFile, &all)
	if err != nil {
		log.Logger().Fatal("Unable to unmarshal all.json", zap.Error(err))
	}

	allMap := make(map[int]All)
	for _, f := range all {
		allMap[f.NeighborhoodID] = f
	}

	uneffectedCount := 0
	// write uneffected to csv file
	unaffectedFile, err := os.Create("generator/uneffected.csv")
	defer unaffectedFile.Close()
	if err != nil {
		log.Logger().Fatal("Unable to create uneffected.csv", zap.Error(err))
	}
	for _, f := range fix {
		if _, ok := allMap[f.NeighborhoodID]; !ok {
			fmt.Println(f.NeighborhoodID)
		}

		if _, ok := allMap[f.NeighborhoodID]; ok {
			allLat, _ := strconv.ParseFloat(allMap[f.NeighborhoodID].Lat, 64)
			allLng, _ := strconv.ParseFloat(allMap[f.NeighborhoodID].Lng, 64)
			if int(f.Lat) == int(allLat) && int(f.Lng) == int(allLng) {
				if uneffectedCount == 0 {
					unaffectedFile.WriteString("nId,gorunen_il,gorunmesi_gereken_il,mahalle\n")
				}
				uneffectedCount++
				unaffectedFile.WriteString(strconv.Itoa(f.NeighborhoodID) + "," + allMap[f.NeighborhoodID].CityNameFetched + "," + allMap[f.NeighborhoodID].CityNameOriginal + "," + allMap[f.NeighborhoodID].Neighborhood + "\n")
				fmt.Println(strconv.Itoa(uneffectedCount)+",", "Gorunen il/ilce:", allMap[f.NeighborhoodID].CityNameFetched+",", "Gorunmesi gereken il:", allMap[f.NeighborhoodID].CityNameOriginal+",", "Mahalle:", allMap[f.NeighborhoodID].Neighborhood)
				//fmt.Println("CityID", f.CityID, "DistrictID", f.DistrictID, "NeighborhoodID", f.NeighborhoodID, "Lat", f.Lat, "Lng", f.Lng)
			}
		}
	}

	fmt.Println("Uneffected count", uneffectedCount)

	// get all neighborhood ids
	//rawSql := psql.Select("id", "new_id").From("neighborhood").Where(squirrel.Eq{"id": neighIds}).OrderBy("id ASC")
	//sql, args, err := rawSql.ToSql()
	//if err != nil {
	//	log.Logger().Fatal("Unable to generate sql", zap.Error(err))
	//}
	//
	//rows, err := pool.Query(context.Background(), sql, args...)
	//if err != nil {
	//	log.Logger().Fatal("Unable to query", zap.Error(err))
	//}
	//
	//neighIdMap := make(map[int]int)
	//for rows.Next() {
	//	var id, newId int
	//	err := rows.Scan(&id, &newId)
	//	if err != nil {
	//		log.Logger().Fatal("Unable to scan", zap.Error(err))
	//	}
	//	neighIdMap[id] = newId
	//}
	//
	//for _, f := range fix {
	//	latStr := strconv.FormatFloat(f.Lat, 'f', -1, 64)
	//	lngStr := strconv.FormatFloat(f.Lng, 'f', -1, 64)
	//	rawSql := psql.Update("locations").Set("latitude", latStr).Set("longitude", lngStr).Where(squirrel.Eq{"neighborhood_id": neighIdMap[f.NeighborhoodID]})
	//	sql, args, err := rawSql.ToSql()
	//	if err != nil {
	//		log.Logger().Fatal("Unable to generate sql", zap.Error(err))
	//	}
	//
	//	fmt.Println("sql", sql)
	//	fmt.Println("args", args)
	//
	//	_, err = pool.Exec(context.Background(), sql, args...)
	//	if err != nil {
	//		log.Logger().Fatal("Unable to exec", zap.Error(err))
	//	}
	//
	//	log.Logger().Info("Updated", zap.Int("neighborhood_id", f.NeighborhoodID), zap.Float64("lat", f.Lat), zap.Float64("lng", f.Lng))
	//}
}

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
			}
		}

		if lastCity.Id != int(location.CityID) {
			response[int64(lastCity.Id)].Districts = cityDistricts[int64(lastCity.Id)]
			response[int64(lastCity.Id)].Lat = lastDistrict.Lat
			response[int64(lastCity.Id)].Lng = lastDistrict.Lng
			cityDistricts[location.CityID] = make([]model.District, 0)
			cityDistrictNeighborhoods[location.CityID] = make(map[int64][]model.Neighborhood)
		}

		if lastDistrict.Id != int(location.DistrictID) {
			lastDistrict.Neighborhoods = cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)]
			var mainNeighborhood *model.Neighborhood
			for _, neighborhood := range lastDistrict.Neighborhoods {
				if strings.Contains(neighborhood.Name, "MERKEZ") {
					mainNeighborhood = &neighborhood
					break
				}
			}
			if mainNeighborhood != nil {
				lastDistrict.Lat = mainNeighborhood.Lat
				lastDistrict.Lng = mainNeighborhood.Lng
			} else if len(lastDistrict.Neighborhoods) > 0 {
				lastDistrict.Lat = lastDistrict.Neighborhoods[0].Lat
				lastDistrict.Lng = lastDistrict.Neighborhoods[0].Lng
			}

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
			DistrictID: int(location.DistrictID),
			CityID:     int(location.CityID),
			Lat:        Lat,
			Lng:        Long,
		}
		cityDistrictNeighborhoods[location.CityID][location.DistrictID] = append(cityDistrictNeighborhoods[location.CityID][location.DistrictID], *neighborhood)

		lastCity = city
		lastDistrict = &model.District{
			Id:     int(location.DistrictID),
			Name:   location.DistrictName,
			CityID: int(location.CityID),
			Lat:    Lat,
			Lng:    Long,
		}
	}

	lastDistrict.Neighborhoods = cityDistrictNeighborhoods[int64(lastCity.Id)][int64(lastDistrict.Id)]
	var mainNeighborhood *model.Neighborhood
	for _, neighborhood := range lastDistrict.Neighborhoods {
		if strings.Contains(neighborhood.Name, "MERKEZ") {
			mainNeighborhood = &neighborhood
			break
		}
	}
	if mainNeighborhood != nil {
		lastDistrict.Lat = mainNeighborhood.Lat
		lastDistrict.Lng = mainNeighborhood.Lng
	} else if len(lastDistrict.Neighborhoods) > 0 {
		lastDistrict.Lat = lastDistrict.Neighborhoods[0].Lat
		lastDistrict.Lng = lastDistrict.Neighborhoods[0].Lng
	}

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
