package feeds

type Feed struct {
	DistrictId    int64 `json:"district_id,omitempty"`
	VolunteerData int   `json:"volunteer_data,omitempty"`
}

type Response struct {
	Count   int    `json:"count"`
	Results []Feed `json:"results"`
}
