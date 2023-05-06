package pgrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	"github.com/xyedo/blindate/pkg/domain/user"
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
	userEntities "github.com/xyedo/blindate/pkg/domain/user/entities"
)

func New(db *sqlx.DB) user.Repository {
	return &userDb{
		conn: db,
	}
}

type userDb struct {
	conn *sqlx.DB
}

func (u *userDb) InsertUser(user userDTOs.RegisterUser) (string, error) {
	args := []any{user.FullName, user.Alias, user.Email, user.Password, user.Dob, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId string
	err := u.conn.GetContext(ctx, &userId, insertUser, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return "", apperror.Conflicted(apperror.Payload{
					Error:   err,
					Message: "email already taken",
				})
			}
			return "", err
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return "", apperror.Timeout(apperror.Payload{Error: err})
		}
		return "", err
	}

	return userId, nil
}

func (u *userDb) UpdateUser(user userEntities.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{user.FullName, user.Alias, user.Email, user.Password, user.Dob, user.Active, time.Now(), user.ID}
	var returnedId string
	err := u.conn.GetContext(ctx, &returnedId, updateUserById, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{
				Error:   err,
				Message: "user not found",
			})
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return apperror.UnprocessableEntity(apperror.PayloadMap{
					Error: err, ErrorMap: map[string]string{
						"email": "already taken",
					},
				})
			}
		}

		return err
	}

	return nil
}

func (u *userDb) GetUserById(id string) (userEntities.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userEntities.User
	err := u.conn.GetContext(ctx, &user, getUserById, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntities.User{}, apperror.NotFound(apperror.Payload{Error: err, Message: "user not found"})
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return userEntities.User{}, apperror.Timeout(apperror.Payload{Error: err})
		}
		return userEntities.User{}, err
	}
	return user, nil
}
func (u *userDb) GetUserByEmail(email string) (userEntities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userEntities.User
	err := u.conn.GetContext(ctx, &user, getUserByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntities.User{}, apperror.NotFound(apperror.Payload{Error: err, Message: "user not found"})
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return userEntities.User{}, apperror.Timeout(apperror.Payload{Error: err})
		}
		return userEntities.User{}, err
	}
	return user, nil
}

func (u *userDb) CreateProfilePicture(userId, pictureRef string, selected bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedId string
	err := u.conn.GetContext(ctx, &returnedId, insertProfilePicture, userId, selected, pictureRef)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", apperror.Timeout(apperror.Payload{Error: err})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", apperror.NotFound(apperror.Payload{Error: err, Message: "user not found"})
			}
			return "", err
		}
		return "", err
	}
	return returnedId, nil
}
func (u *userDb) SelectProfilePicture(userId string, params userDTOs.ProfilePictureQuery) ([]userEntities.UserProfilePic, error) {
	query := selectProfilePicture
	args := []any{userId}

	if params.Selected != nil {
		query += ` AND selected = $2`
		args = append(args, *params.Selected)
	}

	query += ` ORDER BY id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profilePictures []userEntities.UserProfilePic
	err := u.conn.SelectContext(ctx, &profilePictures, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(apperror.Payload{Error: err, Message: "profile picture not found"})
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, apperror.Timeout(apperror.Payload{Error: err})
		}
		return nil, err
	}
	return profilePictures, nil
}

func (u *userDb) UpdateProfilePictureToFalse(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedProfilePictureId string
	err := u.conn.GetContext(ctx, &returnedProfilePictureId, updateProfilePictureToFalse, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{Error: err})
		}

		return err
	}

	return nil

}
