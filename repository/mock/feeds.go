package mock

import (
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"math/rand"
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
