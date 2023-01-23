package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/user"
	userEntity "github.com/xyedo/blindate/pkg/domain/user/entities"
)

func NewUser(db *sqlx.DB) *UserCon {
	return &UserCon{
		conn: db,
	}
}

type UserCon struct {
	conn *sqlx.DB
}

func (u *UserCon) InsertUser(user userEntity.Register) (string, error) {
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
	args := []any{user.FullName, user.Alias, user.Email, user.Password, user.Dob, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userId string
	err := u.conn.GetContext(ctx, &userId, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return "", common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "email already taken")
			}
			return "", pqErr
		}
		if errors.Is(err, context.Canceled) {
			return "", common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return "", err
	}

	return userId, nil
}

func (u *UserCon) UpdateUser(user userEntity.FullDTO) error {
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
	args := []any{user.FullName, user.Alias, user.Email, user.Password, user.Dob, user.Active, time.Now(), user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var retId string
	err := u.conn.GetContext(ctx, &retId, query, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrResourceNotFound)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "email already taken")
			}
		}
		return err
	}
	return nil
}

func (u *UserCon) GetUserById(id string) (userEntity.FullDTO, error) {
	query := `
		SELECT 
			id, alias, full_name, email, "password",active, dob, created_at, updated_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userEntity.FullDTO
	err := u.conn.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity.FullDTO{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return userEntity.FullDTO{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return userEntity.FullDTO{}, err
	}
	return user, nil
}
func (u *UserCon) GetUserByEmail(email string) (userEntity.FullDTO, error) {
	query := `
		SELECT 
			id, email, "password"
		FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userEntity.FullDTO
	err := u.conn.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity.FullDTO{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return userEntity.FullDTO{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return userEntity.FullDTO{}, err
	}
	return user, nil
}

func (u *UserCon) CreateProfilePicture(userId, pictureRef string, selected bool) (string, error) {
	query := `
	INSERT INTO profile_picture(user_id,selected,picture_ref)
	VALUES($1,$2,$3) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := u.conn.GetContext(ctx, &id, query, userId, selected, pictureRef)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "profile picture not found")
			}
			return "", pqErr
		}
		return "", err
	}
	return id, nil
}
func (u *UserCon) SelectProfilePicture(userId string, params *user.ProfilePicQuery) ([]userEntity.ProfilePic, error) {
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
	var profilePics []userEntity.ProfilePic
	err := u.conn.SelectContext(ctx, &profilePics, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.WrapError(err, common.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return nil, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return nil, err
	}
	return profilePics, nil
}

func (u *UserCon) ProfilePicSelectedToFalse(userId string) (int64, error) {
	query := `
	UPDATE profile_picture SET
		selected = false
	WHERE user_id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := u.conn.ExecContext(ctx, query, userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return 0, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row, nil

}

// func (u *userCon) DeleteProfilePic(userId string, id)
