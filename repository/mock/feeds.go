package mock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/repository"
)

func GetFeeds() (*feeds.Response, error) {
	//var response []feeds.Feed

	return release()

	//for _, nh := range repository.NeighborhoodIdToMap {
	//	response = append(response, feeds.Feed{
	//		NeighborhoodId: nh.Id,
	//		VolunteerData:  rand.Intn(5) + 1,
	//	})
	//}
	//
	//return &feeds.Response{
	//	Count:   len(response),
	//	Results: response,
	//}, nil
}

func GetFeedDetail(neighborhoodId int) (*feeds.FeedDetailResponse, error) {
	var response feeds.FeedDetailResponse
	for _, nh := range repository.NeighborhoodIdToMap {
		if nh.Id == neighborhoodId {
			response.NeighborhoodId = nh.Id
			response.Intensity = rand.Intn(5) + 1
			response.LastUpdateTime = time.Now().Add(-time.Minute * time.Duration(rand.Intn(60))).Format(time.RFC3339)
			response.Details = make([]feeds.FeedDetail, 0)
			for i := 0; i < rand.Intn(10)+1; i++ {
				response.Details = append(response.Details, feeds.FeedDetail{
					BuildingName: fmt.Sprintf("Building %d", i),
					BallotBoxNos: []int{
						i,
						i + 1,
						i + 2,
					},
				})
			}
			break
		}
	}
	return &response, nil
}

func release() (*feeds.Response, error) {
	response := make([]feeds.Feed, 0)
	ovoBuildingStore := repository.OvoBuildingStore
	for _, building := range ovoBuildingStore.BuildingInfos {
		nId := repository.CityToDistrictToNeighborhoodToNeighborhoodId[building.City][building.District][building.Neighborhood]
		if nId == 0 {
			// log
			//fmt.Println("Neighborhood not found for building: ",
			//	"city: ", building.City,
			//	"district: ", building.District,
			//	"neighborhood: ", building.Neighborhood,
			//)

			continue
		}

		response = append(response, feeds.Feed{
			NeighborhoodId: nId,
			VolunteerData:  ovoBuildingStore.NeighToAvgScore[building.Neighborhood],
		})
	}

	return &feeds.Response{
		Count:   len(response),
		Results: response,
	}, nil
}
