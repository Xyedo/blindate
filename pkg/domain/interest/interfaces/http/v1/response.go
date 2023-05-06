package v1

import "time"

type bio struct {
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type getInterestDetailResponse struct {
	bio
	Hobbies     []hobbie     `json:"hobbies"`
	MovieSeries []movieSerie `json:"movie_series"`
	Travels     []travel     `json:"travels"`
	Sports      []sport      `json:"sports"`
}
