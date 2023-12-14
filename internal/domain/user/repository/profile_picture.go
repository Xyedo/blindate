package repository

import (
	"context"
	"errors"

	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func InsertProfilePicture(ctx context.Context, conn pg.Querier, profilePicture entities.ProfilePicture) (string, error) {
	const insertProfilePicture = `
	INSERT INTO profile_picture(id, account_id, selected, file_id)
	VALUES($1,$2,$3,$4)`

	var returnedId string
	err := conn.QueryRow(ctx, insertProfilePicture,
		profilePicture.Id,
		profilePicture.UserId,
		profilePicture.Selected,
		profilePicture.FileId,
	).Scan(&returnedId)
	if err != nil {
		return "", err
	}

	return returnedId, nil
}

func UpdateProfilePictureSelectedToFalseByUserId(ctx context.Context, conn pg.Querier, id string) error {
	const query = `
	UPDATE profile_pictures SET
	selected = false
	WHERE account_id = $1
	`

	tag, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("invalid")
	}

	return nil
}
