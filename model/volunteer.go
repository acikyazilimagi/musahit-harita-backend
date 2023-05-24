package model

type VolunteerDoc struct {
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	KvkkAccepted   bool   `json:"kvkkAccepted"`
	NeighborhoodId int    `json:"neighborhoodId"`
}

type Volunteer struct {
	Id           int    `json:"id"`
	VolunteerDoc []byte `json:"volunteerDoc"`
	BuildingId   int    `json:"buildingId,omitempty"`
	LocationId   int    `json:"locationId,omitempty"`
	Confirmed    bool   `json:"confirmed"`
	SourceId     int8   `json:"sourceId"`
}
