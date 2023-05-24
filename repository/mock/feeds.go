package mock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/repository"
)

func GetFeeds() (*feeds.Response, error) {
	var response []feeds.Feed

	for _, nh := range repository.NeighborhoodIdToMap {
		response = append(response, feeds.Feed{
			NeighborhoodId: nh.Id,
			VolunteerData:  rand.Intn(5) + 1,
		})
	}

	return &feeds.Response{
		Count:   len(response),
		Results: response,
	}, nil
}

func GetFeedDetail(neighborhoodId int) (*feeds.FeedDetailResponse, error) {
	var response feeds.FeedDetailResponse
	for _, nh := range repository.NeighborhoodIdToMap {
		if nh.Id == neighborhoodId {
			response.NeighborhoodId = nh.Id
			response.Intensity = rand.Intn(5) + 1
			response.LastUpdateTime = time.Now().Add(-time.Minute * time.Duration(rand.Intn(60))).Format(time.RFC3339)
			response.Details = make([]string, 0)
			for i := 0; i < rand.Intn(10)+1; i++ {
				str := fmt.Sprintf("%s - %d", nh.Name, rand.Intn(1000)+1)
				response.Details = append(response.Details, str)
			}
			break
		}
	}
	return &response, nil
}
