package model

type Neighborhood struct {
	Id         int    `json:"id"`
	DistrictID int    `json:"districtID"`
	CityID     int    `json:"cityID"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Geo        *Geo   `json:"geo"`
}

type Geo struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
