package entities

import (
	"errors"
	"io"
	"net/http"

	apperror "github.com/xyedo/blindate/internal/common/app-error"
)

func (u UserDetail) ToFileIds() ([]string, map[string]int) {
	ids := make([]string, 0, len(u.ProfilePictures))
	fileIdIdx := make(map[string]int, len(u.ProfilePictures))
	for i := range u.ProfilePictures {
		ids = append(ids, u.ProfilePictures[i].FileId)
		fileIdIdx[u.ProfilePictures[i].FileId] = i
	}
	return ids, fileIdIdx
}
func SanitizeMimeType(file io.ReadSeeker, validMimeTypes []string) (string, error) {
	buff := make([]byte, 512)
	bytesRead, err := file.Read(buff)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buff[:bytesRead])

	var isValidType bool
	for _, validMimeType := range validMimeTypes {
		if contentType == validMimeType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return "", apperror.BadPayloadWithPayloadMap(apperror.PayloadMap{
			Payloads: []apperror.ErrorPayload{
				{
					Code: PhotoInvalidType,
					Details: map[string][]string{
						"file": {"invalid MIME-type"},
					},
				},
			},
		})
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return contentType, nil
}
