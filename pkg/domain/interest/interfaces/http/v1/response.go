package v1

import "time"

type bio struct {
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type hobbie struct {
	Id     string `json:"id"`
	Hobbie string `json:"hobbie"`
}

type movieSerie struct {
	Id         string `json:"id"`
	MovieSerie string `json:"movie_serie"`
}

type travel struct {
	Id     string `json:"id"`
	Travel string `json:"travel"`
}

type sport struct {
	Id    string `json:"id"`
	Sport string `json:"sport"`
}

type getInterestDetailResponse struct {
	bio
	Hobbies     []hobbie     `json:"hobbies"`
	MovieSeries []movieSerie `json:"movie_series"`
	Travels     []travel     `json:"travels"`
	Sports      []sport      `json:"sports"`
}
