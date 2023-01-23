package interestEntity

// HobbieDTO one to many with bio
type HobbieDTO struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Hobbie string `json:"hobbie" db:"hobbie" binding:"required,min=2,max=50"`
}

// MovieSerieDTO one to many with bio
type MovieSerieDTO struct {
	Id         string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	MovieSerie string `json:"movieSerie" db:"movie_serie" binding:"required,min=2,max=50"`
}

// TravelDTO one to many with bio
type TravelDTO struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Travel string `json:"travel" db:"travel" binding:"required,min=2,max=50"`
}

// SportDTO one to many with bio
type SportDTO struct {
	Id    string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Sport string `json:"sport" db:"sport" binding:"required,min=2,max=50"`
}

// FullDTO one to one with user
type FullDTO struct {
	BioDTO
	Hobbies     []HobbieDTO     `json:"hobbies"`
	MovieSeries []MovieSerieDTO `json:"movieSeries"`
	Travels     []TravelDTO     `json:"travels"`
	Sports      []SportDTO      `json:"sports"`
}
