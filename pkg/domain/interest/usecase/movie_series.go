package usecase

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateMovieSeriesByInterestId implements interest.Usecase
func (i *interestUC) CreateMovieSeriesByInterestId(
	id string,
	movieSeries []string,
) ([]string, error) {
	movieSeriesDB := make(
		[]interestEntities.MovieSerie,
		0,
		len(movieSeries),
	)
	for _, movieSerie := range movieSeries {
		movieSeriesDB = append(
			movieSeriesDB,
			interestEntities.MovieSerie{
				MovieSerie: movieSerie,
			})
	}

	err := i.interestRepo.CheckInsertMovieSeriesValid(
		id,
		len(movieSeries),
	)
	if err != nil {
		return nil, err
	}

	err = i.interestRepo.InsertMovieSeriesByInterestId(
		id,
		movieSeriesDB,
	)
	if err != nil {
		return nil, err
	}

	returnedIds := make([]string, 0, len(movieSeriesDB))
	for _, movieSerieDB := range movieSeriesDB {
		returnedIds = append(returnedIds, movieSerieDB.Id)
	}

	return returnedIds, nil
}

// UpdateMovieSeriesByInterestId implements interest.Usecase
func (i *interestUC) UpdateMovieSeries(
	movieSeries []interestDTOs.MovieSerie,
) error {
	movieSeriesEntity := make(
		[]interestEntities.MovieSerie,
		0,
		len(movieSeries),
	)
	for _, movieSerie := range movieSeries {
		movieSeriesEntity = append(
			movieSeriesEntity,
			interestEntities.MovieSerie(movieSerie),
		)
	}

	err := i.interestRepo.UpdateMovieSeries(movieSeriesEntity)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMovieSeriesByInterestId implements interest.Usecase
func (i *interestUC) DeleteMovieSeriesByIDs(
	movieSerieIds []string,
) error {
	return i.interestRepo.DeleteMovieSeriesByIDs(movieSerieIds)
}
