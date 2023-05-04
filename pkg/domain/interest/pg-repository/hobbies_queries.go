package repository

const insertHobbies = `
	INSERT INTO hobbies 
		(interest_id, id, hobbie)
	VALUES %s RETURNING id`

const getHobbies = `
	SELECT 
		id,
		hobbie
	FROM hobbies
	WHERE interest_id = $1`

const updateHobbies = `
	UPDATE hobbies AS h SET
		hobbies = h2.hobbie
		FROM (
			VALUES %s
		) AS h2 (id, hobbie)
		WHERE h2.id = h.id
	RETURNING id
	`

const deleteHobbies = `
	DELETE FROM hobbies
	WHERE 
		id IN (%s)
	RETURNING id`

const hobbiesStatistic = `
	UPDATE interest_statistics SET
		hobbie_count = hobbie_count + $1
	WHERE interest_id = $2
	RETURNING interest_id
`
