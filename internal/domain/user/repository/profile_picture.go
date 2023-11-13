package repository

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func InsertProfilePicture(ctx context.Context, conn pg.Querier, profilePicture entities.ProfilePicture) (string, error) {
	const insertProfilePicture = `
	INSERT INTO profile_pictures(id, account_id, selected, file_id)
	VALUES($1,$2,$3,$4)
	RETURNING id`

	var returnedUUID string
	err := conn.QueryRow(ctx, insertProfilePicture,
		profilePicture.UUID,
		profilePicture.UserId,
		profilePicture.Selected,
		profilePicture.FileId,
	).Scan(&returnedUUID)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func UpdateProfilePictureSelectedToFalseByUserId(ctx context.Context, conn pg.Querier, id string) error {
	const query = `
	UPDATE profile_pictures SET
	selected = false
	WHERE account_id = $1
	`

	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
