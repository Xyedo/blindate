package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"golang.org/x/crypto/bcrypt"
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
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)

	query := `
	INSERT INTO users(full_name, email, "password", dob, created_at, updated_at)
	VALUES($1,$2,$3,$4,$5,$5) RETURNING id`
	args := []any{user.FullName, user.Email, user.Password, user.Dob, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = u.GetContext(ctx, &user.ID, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (u *userCon) GetUserByEmail(email string) (*domain.User, error) {
	query := `
		SELECT 
			id, email, "password"
		FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := u.DB.QueryRowxContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
