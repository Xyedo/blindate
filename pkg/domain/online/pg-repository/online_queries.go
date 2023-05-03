package repository

const insertOnline = `
	INSERT INTO onlines (
		user_id,
		last_online,
		is_online
	)
	VALUES ($1,$2,$3)
	RETURNING user_id`

const selectOnline = `	
	SELECT
		user_id, 
		last_online, 
		is_online
	FROM onlines
	WHERE user_id = $1`

const updateOnline = `
	UPDATE onlines SET
		is_online=$1, 
		last_online=COALESCE($2, last_online)
	WHERE user_id=$3
	RETURNING user_id`
