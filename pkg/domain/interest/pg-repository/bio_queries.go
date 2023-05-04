package repository

const insertBio = `
	INSERT INTO interests (
		user_id, 
		bio, 
		created_at, 
		updated_at)
	VALUES($1,$2,$3,$3) 
	RETURNING id`

const getBioByUserId = `
	SELECT
		id, 
		user_id, 
		bio, 
		created_at, 
		updated_at 
	FROM interests
	WHERE user_id = $1`

const updateInterestBio = `
	UPDATE interests SET 
		bio = $1, 
		updated_at = $2  
		WHERE user_id = $3 
	RETURNING id`
