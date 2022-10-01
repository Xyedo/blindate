package entity

import "time"

type Interest struct {
	Id             string    `db:"id"`
	UserId         string    `db:"user_id"`
	Hobbies        []string  `db:"hobbies"`
	MoviesSeries   []string  `db:"movie_series"`
	Traveling      []string  `db:"traveling"`
	Sport          []string  `db:"sport"`
	Bio            string    `db:"bio"`
	SpotifyConnect string    `db:"spotify_connect"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
