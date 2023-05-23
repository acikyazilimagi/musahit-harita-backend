package model

type District struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	CityID int    `json:"cityID"`
	Geo    struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"geo"`
	Neighborhoods []Neighborhood `json:"neighborhoods"`
}
