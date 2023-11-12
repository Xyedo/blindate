package entities

import "time"

const (
	FileTypePhotoProfile = "PHOTO_PROFILE"
)

type File struct {
	UUID      string
	FileType  string
	BlobLink  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}
