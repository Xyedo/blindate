package entities

import (
	"io"

	"github.com/gabriel-vasile/mimetype"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
)

func (u UserDetail) ToFileIds() ([]string, map[string]int) {
	ids := make([]string, 0, len(u.ProfilePictures))
	fileIdToIdx := make(map[string]int, len(u.ProfilePictures))
	for i := range u.ProfilePictures {
		ids = append(ids, u.ProfilePictures[i].FileId)
		fileIdToIdx[u.ProfilePictures[i].FileId] = i
	}
	return ids, fileIdToIdx
}

func (users UserDetails) ToFileIds() ([]string, map[string][2]int) {
	ids := make([]string, 0, len(users))
	fileIdToIdx := make(map[string][2]int, len(users))

	for i, user := range users {
		for j, profilePic := range user.ProfilePictures {
			ids = append(ids, profilePic.FileId)
			fileIdToIdx[profilePic.FileId] = [2]int{i, j}
		}
	}

	return ids, fileIdToIdx
}

func SanitizeMimeType(file io.ReadSeeker, validMimeTypes []string) (*mimetype.MIME, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	contentType, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, err
	}

	var isValidType bool
	for _, validMimeType := range validMimeTypes {
		if contentType.String() == validMimeType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return nil, apperror.BadPayloadWithPayloadMap(apperror.PayloadMap{
			Payloads: []apperror.ErrorPayload{
				{
					Code: ErrCodePhotoInvalidType,
					Details: map[string][]string{
						"file": {"invalid MIME-type"},
					},
				},
			},
		})
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return contentType, nil
}
