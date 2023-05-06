package repository

const insertTravels = `
	INSERT INTO traveling 
		(interest_id, id, travel)
	VALUES %s RETURNING id`

const checkInsertTravelsValid = `
	SELECT count(id)
	FROM traveling
	WHERE interest_id = $1`

const getTravels = `
	SELECT 
		id,
		travel
	FROM traveling
	WHERE interest_id = $1`

const updateTravels = `
	UPDATE traveling AS t SET
	travel = t2.travel
		FROM (VALUES 
			%s
		) AS t2 (id, travel)
		WHERE t2.id = t.id
	RETURNING t.id
	`

const deleteTravels = `
	DELETE FROM traveling
	WHERE 
		id IN (%s)
	RETURNING id`
