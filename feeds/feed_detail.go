package feeds

import "time"

type FeedDetail struct {
	BuildingId   int    `json:"buildingId,omitempty"`
	BuildingName string `json:"buildingName,omitempty"`
	BallotBoxId  int    `json:"ballotBoxId,omitempty"`
}

type FeedDetailResponse struct {
	NeighborhoodId int          `json:"neighborhoodId"`
	LastUpdateTime *time.Time   `json:"lastUpdateTime,omitempty"`
	Intensity      int          `json:"intensity"`
	Details        []FeedDetail `json:"details"`
}
