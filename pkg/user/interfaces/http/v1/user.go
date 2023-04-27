package v1

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/attachment"
	attachmentDTOs "github.com/xyedo/blindate/pkg/attachment/dtos"
	"github.com/xyedo/blindate/pkg/common/app-error/httperror"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/user"
	userDTOs "github.com/xyedo/blindate/pkg/user/dtos"
)

func NewUserHandler(userUC user.Usecase, attachment attachment.Repository) *userH {
	return &userH{
		userUC:     userUC,
		attachment: attachment,
	}
}

type userH struct {
	userUC     user.Usecase
	attachment attachment.Repository
}

func (u *userH) postUserHandler(c *gin.Context) {
	var request userDTOs.RegisterUser

	if err := c.ShouldBindJSON(&request); err != nil {
		httperror.HandleError(c, err)
		return
	}

	userId, err := u.userUC.CreateUser(request)
	if err != nil {
		httperror.HandleError(c, err)
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

func (u *userH) putUserImageProfileHandler(c *gin.Context) {
	selectedQ := c.Query("selected")
	selected := strings.EqualFold(selectedQ, "true")
	userId := c.GetString(constant.KeyUserId)
	key := u.uploadProfilePic(c)
	if key == "" {
		return
	}

	newProfPic := userDTOs.RegisterProfilePicture{
		UserId:      userId,
		Selected:    selected,
		PictureLink: key,
	}
	id, err := u.userUC.CreateNewProfilePic(newProfPic)
	if err != nil {
		httperror.HandleError(c, err)
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
func (u *userH) getUserByIdHandler(c *gin.Context) {
	userId := c.GetString(constant.KeyUserId)

	user, err := u.userUC.GetUserById(userDTOs.GetUserDetail{
		Id: userId,
	})
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func (u *userH) patchUserByIdHandler(c *gin.Context) {
	var request userDTOs.UpdateUser
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	userId := c.GetString(constant.KeyUserId)
	request.Id = userId

	err = u.userUC.UpdateUser(request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated",
	})

}

func (u *userH) uploadProfilePic(c *gin.Context) string {
	var validImageTypes = []string{
		"image/avif",
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/svg+xml",
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "content-Type header is not valid"})
			return ""
		}
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "request did not contain a file"})
			return ""
		}
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"status":  "fail",
				"message": "max byte to upload is 8mB",
			})
			return ""
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return ""
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return ""
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return ""
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return ""
	}

	var isValidType bool
	contentType := http.DetectContentType(buff)
	for _, validTypes := range validImageTypes {
		if contentType == validTypes {
			isValidType = true
			break
		}
	}
	if !isValidType {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "not valid mime-type"})
		return ""
	}
	key, err := u.attachment.UploadBlob(file, attachmentDTOs.Uploader{
		Length:      fileHeader.Size,
		ContentType: contentType,
		Prefix:      "profile-pictures",
		Ext:         filepath.Ext(fileHeader.Filename),
	})
	if err != nil {
		httperror.HandleError(c, err)
		return ""
	}
	return key
}
