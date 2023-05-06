package repository

const insertMovieSeries = `
	INSERT INTO movie_series 
		(interest_id, id, movie_serie)
	VALUES %s RETURNING id`

const checkInsertMovieSeriesValid = `
	SELECT count(id)
	FROM movie_series
	WHERE interest_id = $1`

const getMovieSeries = `
	SELECT 
		id,
		movie_serie
	FROM movie_series
	WHERE interest_id = $1`

const updateMovieSeries = `
	UPDATE movie_series AS m SET
		movie_serie = m2.movie_serie
		FROM (VALUES 
			%s
		) AS m2 (id, movie_serie)
		WHERE m2.id = m.id
	RETURNING m.id
	`

const deleteMovieSeries = `
	DELETE FROM movie_series
	WHERE 
		id IN (%s)
	RETURNING id`
