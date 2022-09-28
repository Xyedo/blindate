package domain

type Location struct {
	UserId string  `json:"userId"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}
