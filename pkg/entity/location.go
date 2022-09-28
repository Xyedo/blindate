package entity

type Location struct {
	UserId string `db:"user_id"`
	Geog   string `db:"geog"`
}
