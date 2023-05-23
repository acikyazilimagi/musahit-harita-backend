package feeds

type Feed struct {
	NeighborhoodId int `json:"neighborhood_id,omitempty"`
	VolunteerData  int `json:"volunteer_data,omitempty"`
}

type Response struct {
	Count   int    `json:"count"`
	Results []Feed `json:"results"`
}
