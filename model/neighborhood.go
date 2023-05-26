package model

type Neighborhood struct {
	Id         int    `json:"id"`
	DistrictId int    `json:"districtId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Geo        *Geo   `json:"geo"`
}

type Geo struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
