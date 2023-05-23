package mock

import (
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"math/rand"
)

func GetFeeds() (*feeds.Response, error) {
	var response []feeds.Feed

	for _, district := range repository.Districts {
		response = append(response, feeds.Feed{
			DistrictId:    district.Id,
			VolunteerData: rand.Intn(4) + 1,
		})
	}

	return &feeds.Response{
		Count:   len(response),
		Results: response,
	}, nil
}
