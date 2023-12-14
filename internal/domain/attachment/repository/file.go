package repository

import (
	"context"
	"errors"

	"github.com/xyedo/blindate/internal/domain/attachment/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func InsertFile(ctx context.Context, conn pg.Querier, file entities.File) (string, error) {
	const insertFile = `
	INSERT INTO file(id,type,blob_link,created_at,updated_at,version)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING id`

	var returnedUUID string
	err := conn.QueryRow(ctx, insertFile,
		file.Id,
		file.FileType,
		file.BlobLink,
		file.CreatedAt,
		file.UpdatedAt,
		file.Version,
	).Scan(&returnedUUID)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func GetFileByIds(ctx context.Context, conn pg.Querier, ids []string) ([]entities.File, error) {
	const getFileById = `
	SELECT 
		id,
		type,
		blob_link,
		created_at,
		updated_at,
		version
	FROM file
	WHERE id IN (?)`

	query, args, err := pg.In(getFileById, ids)
	if err != nil {
		return nil, err
	}

	files := make([]entities.File, 0)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file entities.File
		err = rows.Scan(
			&file.Id,
			&file.FileType,
			&file.BlobLink,
			&file.CreatedAt,
			&file.UpdatedAt,
			&file.Version,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) != len(ids) {
		return nil, errors.New("invalid")
	}

	return files, nil
}
