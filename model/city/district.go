package city

type District struct {
	Id        int64   `json:"id"`
	CityId    int64   `json:"city_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
