package attachment

import (
	"io"

	attachmentDTOs "github.com/xyedo/blindate/pkg/domain/attachment/dtos"
)

type Repository interface {
	UploadBlob(io.Reader, attachmentDTOs.Uploader) (string, error)
	DeleteBlob(key string) error
	GetPresignedUrl(key string) (string, error)
}
