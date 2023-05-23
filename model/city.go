package model

type City struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Districts []District `json:"districts"`
}
