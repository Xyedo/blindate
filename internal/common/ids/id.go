package ids

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func nanoid() string {
	return gonanoid.MustGenerate("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz", 16)
}
func Hobbie() string {
	return "hobbie_" + nanoid()
}
func MovieSerie() string {
	return "movie_" + nanoid()
}

func Travel() string {
	return "travel_" + nanoid()
}

func Sport() string {
	return "sport_" + nanoid()
}

func ProfilePicture() string {
	return "profilepic_" + nanoid()
}

func File() string {
	return "file_" + nanoid()
}

func Attachment() string {
	return "attachment_" + nanoid()
}
