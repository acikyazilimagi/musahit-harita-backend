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
)

var OvoBuildingStore *OvoBuildingsInfo

var ovoBuildingData []byte

// OvoBuilding struct, represents each line of data
type OvoBuilding struct {
	City          string
	District      string
	DistrictScore int
	Neighborhood  string
	School        string
	SchoolScore   int
	SchoolID      string
}

type OvoBuildingsInfo struct {
	BuildingInfos          []OvoBuilding
	CityToDistrictsToNeigh map[string]map[string]map[string][]OvoBuilding
	NeighToAvgScore        map[string]int
}

func NewOvoBuildingInfo(data s3.ObjectData) *OvoBuildingsInfo {
	ovoBuildingData = data.Bytes()
	return &OvoBuildingsInfo{
		BuildingInfos:          make([]OvoBuilding, 0),
		CityToDistrictsToNeigh: make(map[string]map[string]map[string][]OvoBuilding),
		NeighToAvgScore:        make(map[string]int),
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
		schoolData := strings.Split(record[2], "|")
		if len(schoolData) != 2 {
			log.Logger().Warn("Invalid school data", zap.Int("line", count), zap.Any("record", record))
			continue
		}
		// split the school name column by "-"
		schoolNameData := strings.Split(schoolData[0], " - ")
		if len(schoolNameData) < 3 {
			log.Logger().Warn("Invalid school name data", zap.Int("line", count), zap.Any("record", record))
			continue
		}

		city := record[0]
		district := stringutils.ParseOvoDistrict(strings.TrimSpace(concatData[1]))
		districtScore, err := strconv.Atoi(strings.TrimSpace(concatData[0]))
		if err != nil {
			log.Logger().Warn("Cannot convert district score to int", zap.Int("line", count), zap.Any("record", record))
			continue
		}

		schoolScore, err := strconv.Atoi(strings.TrimSpace(schoolNameData[0]))
		if err != nil {
			log.Logger().Warn("Cannot convert school score to int", zap.Int("line", count), zap.Any("record", record))
		}

		if len(strings.Split(schoolNameData[1], "-")) > 1 {
			schoolNameData[1] = strings.Split(schoolNameData[1], "-")[0]
		}
		neighborhood := stringutils.ParseOvoNeighborhood(strings.TrimSpace(schoolNameData[1]))
		schoolName := strings.TrimSpace(schoolNameData[2])
		schoolID := schoolData[1]

		ovoBuilding := OvoBuilding{
			City:          city,
			District:      district,
			DistrictScore: districtScore,
			Neighborhood:  neighborhood,
			School:        schoolName,
			SchoolScore:   schoolScore,
			SchoolID:      schoolID,
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
		o.NeighToAvgScore[neighborhood] = int(math.Ceil(float64((o.NeighToAvgScore[neighborhood] + schoolScore) / 2.0)))
	}

	OvoBuildingStore = o
	return OvoBuildingStore
}
