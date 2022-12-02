package api

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type userSvc interface {
	CreateUser(newUser *domain.User) error
	VerifyCredential(email, password string) (string, error)
	GetUserById(id string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	CreateNewProfilePic(profPicParam domain.ProfilePicture) (string, error)
}

type attachmentManager interface {
	UploadBlob(file io.Reader, attach domain.Uploader) (string, error)
}

func NewUser(userSvc userSvc, attachmentSvc attachmentManager) *user {
	return &user{
		userService:   userSvc,
		attachmentSvc: attachmentSvc,
	}
}

type user struct {
	userService   userSvc
	attachmentSvc attachmentManager
}

func (u *user) postUserHandler(c *gin.Context) {
	var input struct {
		FullName string    `json:"fullName" binding:"required,max=50"`
		Alias    string    `json:"alias" binding:"required,max=15"`
		Email    string    `json:"email" binding:"required,email"`
		Password string    `json:"password" binding:"required,min=8"`
		Dob      time.Time `json:"dob" binding:"required,validdob"`
	}

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

	user := domain.User{
		FullName: input.FullName,
		Alias:    input.Alias,
		Email:    input.Email,
		Password: input.Password,
		Dob:      input.Dob,
	}
	err := u.userService.CreateUser(&user)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "confirmation email sent, check your email!",
		"data": map[string]string{
			"id": user.ID,
		},
	})
}

func (u *user) putUserImageProfileHandler(c *gin.Context) {
	selectedQ := c.Query("selected")
	selected := strings.EqualFold(selectedQ, "true")
	userId := c.GetString("userId")
	userDb, err := u.userService.GetUserById(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
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
	newProfPic := domain.ProfilePicture{
		UserId:      userDb.ID,
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
func (u *user) getUserByIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
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

func (u *user) patchUserByIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	user, err := u.userService.GetUserById(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	var input struct {
		FullName    *string    `json:"fullName" binding:"omitempty,max=50"`
		Alias       *string    `json:"alias" binding:"omitempty,max=15"`
		Email       *string    `json:"email" binding:"omitempty,email"`
		OldPassword *string    `json:"oldPassword" binding:"required_with=NewPassword,omitempty,min=8"`
		NewPassword *string    `json:"newPassword" binding:"required_with=OldPassword,omitempty,min=8"`
		Dob         *time.Time `json:"dob" binding:"omitempty,validdob"`
	}

	err = c.ShouldBindJSON(&input)
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
	if input.NewPassword != nil && input.OldPassword != nil {
		_, err := u.userService.VerifyCredential(user.Email, *input.OldPassword)
		if err != nil {
			jsonHandleError(c, err)
			return
		}
		user.Password = *input.NewPassword
	}

	if input.FullName != nil {
		user.FullName = *input.FullName
	}

	if input.Alias != nil {
		user.Alias = *input.Alias
	}

	if input.Email != nil {
		user.Active = false
		user.Email = *input.Email
	}
	if input.Dob != nil {
		user.Dob = *input.Dob
	}

	err = u.userService.UpdateUser(user)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated",
	})

}
