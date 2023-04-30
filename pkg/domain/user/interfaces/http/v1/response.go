package v1

import (
	"time"
)

type getUserResponse struct {
	ID         string               `json:"id"`
	FullName   string               `json:"full_name"`
	Alias      string               `json:"alias"`
	ProfilePic []userProfilePicture `json:"profile_picture"`
	Email      string               `json:"email"`
	Dob        time.Time            `json:"dob"`
}

type userProfilePicture struct {
	Id          string `json:"id"`
	Selected    bool   `json:"selected"`
	PictureLink string `json:"picture_ref"`
}
