package domain

type Location struct {
	UserId string `json:"-"`
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
}
