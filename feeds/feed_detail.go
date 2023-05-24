package feeds

type FeedDetail struct {
	BuildingName     string
	BallotBoxCombine string
}

type FeedDetailResponse struct {
	NeighborhoodId int      `json:"neighborhoodId"`
	LastUpdateTime string   `json:"lastUpdateTime,omitempty"`
	Intensity      int      `json:"intensity"`
	Details        []string `json:"details"`
}
