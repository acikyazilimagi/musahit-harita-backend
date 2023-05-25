package repository

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"github.com/acikkaynak/musahit-harita-backend/aws/s3"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/acikkaynak/musahit-harita-backend/utils/stringutils"
	"go.uber.org/zap"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var OvoBuildingStore *OvoBuildingsInfo

var ovoBuildingData []byte

// OvoBuilding struct, represents each line of data
type OvoBuilding struct {
	City          string
	District      string
	DistrictScore int
	Neighborhood  string
	BuildingName  string
	BuildingScore int
	BuildingId    string
}

type OvoBuildingsInfo struct {
	Mutex                     sync.RWMutex
	BuildingInfos             []OvoBuilding
	CityToDistrictsToNeigh    map[string]map[string]map[string][]OvoBuilding
	NeighToAvgScore           map[string]int
	NeighToId                 map[string]int
	NeighborhoodIdToBuildings map[int][]OvoBuilding
	NeighborhoodIdToAvgScore  map[int]int
	LastUpdateTime            time.Time
}

func NewOvoBuildingInfo(data s3.ObjectData) *OvoBuildingsInfo {
	ovoBuildingData = data.Bytes()
	return &OvoBuildingsInfo{
		BuildingInfos:             make([]OvoBuilding, 0),
		CityToDistrictsToNeigh:    make(map[string]map[string]map[string][]OvoBuilding),
		NeighToAvgScore:           make(map[string]int),
		NeighToId:                 make(map[string]int),
		NeighborhoodIdToBuildings: make(map[int][]OvoBuilding),
		NeighborhoodIdToAvgScore:  make(map[int]int),
		LastUpdateTime:            time.Now(),
	}
}

func (o *OvoBuildingsInfo) Store() *OvoBuildingsInfo {
	// ovoBuildingData is a byte array, we need to convert it to a io.Reader
	f := strings.NewReader(string(ovoBuildingData))

	reader := csv.NewReader(f)
	reader.Comma = ',' // set comma as column separator
	reader.ReuseRecord = true

	// Read the file line by line
	count := 0
	for {
		record, err := reader.Read()
		if err != nil || record == nil {
			break
		}
		if count == 0 {
			count++
			continue
		}
		count++
		if len(record) != 3 {
			log.Logger().Warn("Invalid record", zap.String("record", fmt.Sprintf("%v", record)))
			continue
		}

		// split the concat column by "-"
		concatData := strings.Split(record[1], "-")
		if len(concatData) != 3 {
			log.Logger().Warn("Invalid concat data", zap.Int("line", count), zap.Any("record", record))
			continue
		}
		// split the last column by "|"
		buildingData := strings.Split(record[2], "|")
		if len(buildingData) != 2 {
			log.Logger().Warn("Invalid building data", zap.Int("line", count), zap.Any("record", record))
			continue
		}
		// split the building name column by "-"
		buildingNameData := strings.Split(buildingData[0], " - ")
		if len(buildingNameData) < 3 {
			log.Logger().Warn("Invalid building name data", zap.Int("line", count), zap.Any("record", record))
			continue
		}

		city := record[0]
		district := stringutils.ParseOvoDistrict(strings.TrimSpace(concatData[1]))
		districtScore, err := strconv.Atoi(strings.TrimSpace(concatData[0]))
		if err != nil {
			log.Logger().Warn("Cannot convert district score to int", zap.Int("line", count), zap.Any("record", record))
			continue
		}

		buildingScore, err := strconv.Atoi(strings.TrimSpace(buildingNameData[0]))
		if err != nil {
			log.Logger().Warn("Cannot convert building score to int", zap.Int("line", count), zap.Any("record", record))
		}

		if len(strings.Split(buildingNameData[1], "-")) > 1 {
			buildingNameData[1] = strings.Split(buildingNameData[1], "-")[0]
		}
		neighborhood := stringutils.ParseOvoNeighborhood(strings.TrimSpace(buildingNameData[1]))
		buildingName := strings.TrimSpace(buildingNameData[2])
		buildingID := buildingData[1]

		ovoBuilding := OvoBuilding{
			City:          city,
			District:      district,
			DistrictScore: districtScore,
			Neighborhood:  neighborhood,
			BuildingName:  buildingName,
			BuildingScore: buildingScore,
			BuildingId:    buildingID,
		}

		o.BuildingInfos = append(o.BuildingInfos, ovoBuilding)

		// if city is not in the map, add it
		if _, ok := o.CityToDistrictsToNeigh[city]; !ok {
			o.CityToDistrictsToNeigh[city] = make(map[string]map[string][]OvoBuilding)
		}
		// if district is not in the map, add it
		if _, ok := o.CityToDistrictsToNeigh[city][district]; !ok {
			o.CityToDistrictsToNeigh[city][district] = make(map[string][]OvoBuilding)
		}
		// if neighborhood is not in the map, add it
		if _, ok := o.CityToDistrictsToNeigh[city][district][neighborhood]; !ok {
			o.CityToDistrictsToNeigh[city][district][neighborhood] = make([]OvoBuilding, 0)
		}
		// append the building to the map
		o.CityToDistrictsToNeigh[city][district][neighborhood] = append(o.CityToDistrictsToNeigh[city][district][neighborhood], ovoBuilding)

		// if neighborhood is not in the map, add it
		if _, ok := o.NeighToAvgScore[neighborhood]; !ok {
			o.NeighToAvgScore[neighborhood] = 0
		}

		// get the average score of the neighborhood
		o.NeighToAvgScore[neighborhood] = int(math.Ceil(float64((o.NeighToAvgScore[neighborhood] + buildingScore) / 2.0)))

		nId := CityToDistrictToNeighborhoodToNeighborhoodId[city][district][neighborhood]
		if nId != 0 {
			o.NeighToId[neighborhood] = nId
			if _, ok := o.NeighborhoodIdToBuildings[nId]; !ok {
				o.NeighborhoodIdToBuildings[nId] = make([]OvoBuilding, 0)
			}
			o.NeighborhoodIdToBuildings[nId] = append(o.NeighborhoodIdToBuildings[nId], ovoBuilding)
			o.NeighborhoodIdToAvgScore[nId] = o.NeighToAvgScore[neighborhood]
		}
	}

	OvoBuildingStore = o
	return OvoBuildingStore
}
