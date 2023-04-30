package pgrepository

const insertUser = `
	INSERT INTO users(
		full_name, 
		alias, 
		email, 
		"password", 
		dob,
		created_at, 
		updated_at
	)
	VALUES($1,$2,$3,$4,$5,$6,$6) RETURNING id
`

const updateUserById = `
	UPDATE users
	SET 
		full_name = $1, 
		alias=$2, 
		email = $3, 
		"password" = $4, 
		dob=$5, 
		active=$6, 
		updated_at = $7
	WHERE id = $8
	RETURNING id
`

const getUserById = `
	SELECT 
		id, alias, full_name, email, "password",active, dob, created_at, updated_at
	FROM users
	WHERE id = $1
`

const getUserByEmail = `
	SELECT 
		id, email, "password"
	FROM users WHERE email = $1
`
const insertProfilePicture = `	
	INSERT INTO 
		profile_picture(user_id,selected,picture_ref)
	VALUES 
		($1,$2,$3) 
	RETURNING id
`
const selectProfilePicture = `
	SELECT 
		id,
		user_id,
		selected,
		picture_ref 
	FROM profile_picture 
	WHERE user_id =$1`

const updateProfilePictureToFalse = `
	UPDATE profile_picture SET
		selected = false
	WHERE 
		user_id=$1
	RETURNING id
`
