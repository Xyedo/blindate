package repository

const insertSports = `
	INSERT INTO sports 
		(interest_id, id, sport)
	VALUES %s RETURNING id`

const checkInsertSportsValid = `
	SELECT count(id)
	FROM sports
	WHERE interest_id = $1`

const getSports = `
	SELECT 
		id,
		sport
	FROM sports
	WHERE interest_id = $1`

const updateSports = `
	UPDATE sports AS s SET
		sport = s2.sport
		FROM (VALUES 
			%s
		) AS s2 (id, sport)
		WHERE s2.id = s.id
	RETURNING s.id
	`

const deleteSports = `
	DELETE FROM sports
	WHERE 
		id IN (%s)
	RETURNING id`
