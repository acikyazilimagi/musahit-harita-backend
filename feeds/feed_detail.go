package feeds

type FeedDetail struct {
	BuildingName string `json:"buildingName"`
	BallotBoxNos []int  `json:"ballotBoxNos"`
}

type FeedDetailResponse struct {
	NeighborhoodId int          `json:"neighborhoodId"`
	LastUpdateTime string       `json:"lastUpdateTime,omitempty"`
	Intensity      int          `json:"intensity"`
	Details        []FeedDetail `json:"details"`
}
