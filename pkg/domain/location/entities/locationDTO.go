package locationEntity

// Location one to one with user
type DTO struct {
	UserId string `json:"-"`
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
}
