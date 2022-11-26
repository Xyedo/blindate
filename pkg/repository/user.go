package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type User interface {
	InsertUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserById(id string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	CreateProfilePicture(userId, pictureRef string, selected bool) (string, error)
	SelectProfilePicture(userId string, params *entity.ProfilePicQuery) ([]domain.ProfilePicture, error)
	ProfilePicSelectedToFalse(userId string) (int64, error)
}

func NewUser(db *sqlx.DB) *userCon {
	return &userCon{
		conn: db,
	}
}

type userCon struct {
	conn *sqlx.DB
}

func (u *userCon) InsertUser(user *domain.User) error {
	query := `
	INSERT INTO users(
		full_name, 
		alias, 
		email, 
		"password", 
		dob,
		created_at, 
		updated_at
	)
	VALUES($1,$2,$3,$4,$5,$6,$6) RETURNING id`
	args := []any{user.FullName, user.Alias, user.Email, user.HashedPassword, user.Dob, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.conn.GetContext(ctx, &user.ID, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (u *userCon) UpdateUser(user *domain.User) error {
	query := `
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
		RETURNING id`
	args := []any{user.FullName, user.Alias, user.Email, user.HashedPassword, user.Dob, user.Active, time.Now(), user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := u.conn.GetContext(ctx, &id, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (u *userCon) GetUserById(id string) (*domain.User, error) {
	query := `
		SELECT 
			id, alias, full_name, email, "password",active, dob, created_at, updated_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := u.conn.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userCon) GetUserByEmail(email string) (*domain.User, error) {
	query := `
		SELECT 
			id, email, "password"
		FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := u.conn.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userCon) CreateProfilePicture(userId, pictureRef string, selected bool) (string, error) {
	query := `
	INSERT INTO profile_picture(user_id,selected,picture_ref)
	VALUES($1,$2,$3) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := u.conn.GetContext(ctx, &id, query, userId, selected, pictureRef)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (u *userCon) SelectProfilePicture(userId string, params *entity.ProfilePicQuery) ([]domain.ProfilePicture, error) {
	query := `
	SELECT 
		id,
		user_id,
		selected,
		picture_ref 
	FROM profile_picture 
	WHERE user_id =$1`
	args := []any{userId}
	if params != nil {
		if params.Selected != nil {
			query += ` AND selected = $2`
			args = append(args, *params.Selected)
		}

	}
	query += ` ORDER BY id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var profilePics []domain.ProfilePicture
	err := u.conn.SelectContext(ctx, &profilePics, query, args...)
	if err != nil {
		return nil, err
	}
	return profilePics, nil
}

func (u *userCon) ProfilePicSelectedToFalse(userId string) (int64, error) {
	query := `
	UPDATE profile_picture SET
		selected = false
	WHERE user_id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := u.conn.ExecContext(ctx, query, userId)
	if err != nil {
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row, nil

}

// func (u *userCon) DeleteProfilePic(userId string, id)
