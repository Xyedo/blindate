package repository

const insertBasicInfo = `
	INSERT INTO basic_info(
		user_id, 
		gender, 
		from_loc, 
		height, 
		education_level,
		drinking,
		smoking,
		relationship_pref,
		looking_for, 
		zodiac, 
		kids, 
		work, 
		created_at, 
		updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $13)
	RETURNING user_id`

const getBasicInfoByUserId = `
	SELECT
		user_id, 
		gender, 
		from_loc, 
		height, 
		education_level, 
		drinking, 
		smoking, 
		relationship_pref, 
		looking_for, 
		zodiac, 
		kids, 
		work, 
		created_at, 
		updated_at
	FROM basic_info
	WHERE user_id = $1
`
const updateBasicInfo = `
	UPDATE basic_info SET
		gender =$1, 
		from_loc=$2, 
		height=$3, 
		education_level=$4,
		drinking=$5,
		smoking=$6,
		relationship_pref=$7,
		looking_for=$8, 
		zodiac=$9, 
		kids=$10, 
		work=$11, 
		updated_at=$12
	WHERE user_id = $13
	RETURNING user_id
`
