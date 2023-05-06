package repository

const insertHobbies = `
	INSERT INTO hobbies 
		(interest_id, id, hobbie)
	VALUES %s RETURNING id`

const checkInsertHobbiesValid = `
	SELECT count(id)
	FROM hobbies
	WHERE interest_id = $1`

const getHobbies = `
	SELECT 
		id,
		hobbie
	FROM hobbies
	WHERE interest_id = $1`

const updateHobbies = `
	UPDATE hobbies AS h SET
		hobbie = h2.hobbie
		FROM (VALUES 
			%s
		) AS h2 (id, hobbie)
		WHERE h2.id = h.id
	RETURNING h.id
	`

const deleteHobbies = `
	DELETE FROM hobbies
	WHERE 
		id IN (%s)
	RETURNING id`
