package model

type District struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	CityID        int            `json:"cityID"`
	Geo           *Geo           `json:"geo"`
	Neighborhoods []Neighborhood `json:"neighborhoods"`
}
