package model

type City struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Lat       float64    `json:"lat"`
	Lng       float64    `json:"lng"`
	Districts []District `json:"districts"`
}
