package model

type Neighborhood struct {
	Id         int     `json:"id"`
	DistrictID int     `json:"districtID"`
	CityID     int     `json:"cityID"`
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}
