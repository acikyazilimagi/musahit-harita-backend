package model

type ElectionResult struct {
	CityName         string `json:"city_name"`
	DistrictName     string `json:"district_name"`
	NeighborhoodName string `json:"neighborhood_name"`
	SchoolName       string `json:"school_name"`
	BoxNumber        int    `json:"box_number"`
	IsApproved       bool   `json:"is_approved"`
	Count            int    `json:"count"`
}
