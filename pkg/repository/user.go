package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
)

func NewUser(db *sqlx.DB) *userCon {
	return &userCon{
		db,
	}
}

type userCon struct {
	*sqlx.DB
}

func (u *userCon) InsertUser(user *domain.User) error {
	query := `
	INSERT INTO users(
		full_name, 
		alias, 
		image_profile, 
		email, 
		"password", 
		dob,
		created_at, 
		updated_at
	)
	VALUES($1,$2,$3,$4,$5,$6,$7,$7) RETURNING id`
	args := []any{user.FullName, user.Alias, user.ImageProfile, user.Email, user.HashedPassword, user.Dob, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.GetContext(ctx, &user.ID, query, args...)
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
			image_profile= $3,
			email = $4, 
			"password" = $5, 
			dob=$6, 
			active=$7, 
			updated_at = $8
		WHERE id = $9
		RETURNING id`
	args := []any{user.FullName, user.Alias, user.ImageProfile, user.Email, user.HashedPassword, user.Dob, user.Active, time.Now(), user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := u.GetContext(ctx, &id, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (u *userCon) GetUserById(id string) (*domain.User, error) {
	query := `
		SELECT 
			id, alias, full_name, image_profile, email, "password",active, dob, created_at, updated_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := u.DB.GetContext(ctx, &user, query, id)
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
	err := u.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
