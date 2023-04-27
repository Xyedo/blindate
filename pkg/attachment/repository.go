package attachment

import (
	"io"

	"github.com/xyedo/blindate/pkg/attachment/dtos"
)

type Repository interface {
	UploadBlob(io.Reader, dtos.Uploader) (string, error)
	DeleteBlob(key string) error
	GetPresignedUrl(key string) (string, error)
}
