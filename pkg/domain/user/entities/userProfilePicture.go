package userEntity

// ProfilePic one to many with user
type ProfilePic struct {
	Id          string `json:"id" db:"id"`
	UserId      string `json:"userId" db:"user_id"`
	Selected    bool   `json:"selected" db:"selected"`
	PictureLink string `json:"pictureLink" db:"picture_ref"`
}
