package model

type District struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	CityID        int            `json:"cityID"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Neighborhoods []Neighborhood `json:"neighborhoods"`
}
