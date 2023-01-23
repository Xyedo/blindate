package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	attachmentEntity "github.com/xyedo/blindate/pkg/domain/attachment"
	userEntity "github.com/xyedo/blindate/pkg/domain/user/entities"
)

type userSvc interface {
	CreateUser(newUser userEntity.Register) (string, error)
	GetUserById(id string) (userEntity.FullDTO, error)
	UpdateUser(userId string, updateUser userEntity.Update) error
	CreateNewProfilePic(profPicParam userEntity.ProfilePic) (string, error)
}

type attachmentManager interface {
	UploadBlob(file io.Reader, attach attachmentEntity.Uploader) (string, error)
}

func NewUser(userSvc userSvc, attachmentSvc attachmentManager) *User {
	return &User{
		userService:   userSvc,
		attachmentSvc: attachmentSvc,
	}
}

type User struct {
	userService   userSvc
	attachmentSvc attachmentManager
}

func (u *User) postUserHandler(c *gin.Context) {
	var input userEntity.Register

	if err := c.ShouldBindJSON(&input); err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"fullname": "must be required and between 1-50 characters",
			"alias":    "must be required and between 1-15 characters",
			"email":    "must be required and have an valid email",
			"password": "must be required and have more than 8 character",
			"dob":      "must be required and between today and after 1990",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	userId, err := u.userService.CreateUser(input)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "confirmation email sent, check your email!",
		"data": map[string]string{
			"id": userId,
		},
	})
}

func (u *User) putUserImageProfileHandler(c *gin.Context) {
	selectedQ := c.Query("selected")
	selected := strings.EqualFold(selectedQ, "true")
	userId := c.GetString(keyUserId)
	var validImageTypes = []string{
		"image/avif",
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/svg+xml",
	}
	key, _ := uploadFile(c, u.attachmentSvc, validImageTypes, "profile-picture")
	if key == "" {
		return
	}
	newProfPic := userEntity.ProfilePic{
		UserId:      userId,
		Selected:    selected,
		PictureLink: key,
	}
	id, err := u.userService.CreateNewProfilePic(newProfPic)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user profile-picture uploaded",
		"data": gin.H{
			"profilePicture": gin.H{
				"id": id,
			},
		},
	})
}
func (u *User) getUserByIdHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	user, err := u.userService.GetUserById(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func (u *User) patchUserByIdHandler(c *gin.Context) {
	var input userEntity.Update
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"fullName":    "must less than 50 character",
			"alias":       "must less than 15 character",
			"email":       "must be an valid email",
			"oldPassword": "must be more than 8 character and pairing with newPassword",
			"newPassword": "must be more than 8 character and pairing with oldPassword",
			"dob":         "must betwen today and more than 1990",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	userId := c.GetString(keyUserId)
	err = u.userService.UpdateUser(userId, input)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated",
	})

}
