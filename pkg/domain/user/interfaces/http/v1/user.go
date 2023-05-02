package v1

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/domain/attachment"
	attachmentDTOs "github.com/xyedo/blindate/pkg/domain/attachment/dtos"
	"github.com/xyedo/blindate/pkg/domain/user"
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
	"github.com/xyedo/blindate/pkg/infrastructure"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(config infrastructure.Config, userUC user.Usecase, attachment attachment.Repository) *userH {
	return &userH{
		config:     config,
		userUC:     userUC,
		attachment: attachment,
	}
}

type userH struct {
	config     infrastructure.Config
	userUC     user.Usecase
	attachment attachment.Repository
}

func (u *userH) postUserHandler(c *gin.Context) {
	var request postUserRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	userId, err := u.userUC.CreateUser(userDTOs.RegisterUser(request))
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
		Id:             userId,
		ProfilePicture: true,
	})
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	jsonProfilePics := make([]userProfilePicture, 0)
	for _, profilePic := range user.ProfilePic {
		jsonProfilePics = append(jsonProfilePics, userProfilePicture{
			Id:          profilePic.Id,
			Selected:    profilePic.Selected,
			PictureLink: profilePic.PictureLink,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": getUserResponse{
				ID:         user.ID,
				FullName:   user.FullName,
				Alias:      user.Alias,
				ProfilePic: jsonProfilePics,
				Email:      user.Email,
				Dob:        user.Dob,
			},
		},
	})
}

func (u *userH) patchUserByIdHandler(c *gin.Context) {
	var request patchUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	userId := c.GetString(constant.KeyUserId)

	err = u.userUC.UpdateUser(userDTOs.UpdateUser{
		Id:          userId,
		FullName:    request.FullName,
		Alias:       request.Alias,
		Email:       request.Email,
		OldPassword: request.OldPassword,
		NewPassword: request.NewPassword,
		Dob:         request.Dob,
	})
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
		httperror.HandleError(c, err)
		return ""
	}
	file, err := fileHeader.Open()
	if err != nil {
		httperror.HandleError(c, err)
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
		httperror.HandleError(c, err)
		return ""
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		httperror.HandleError(c, err)
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
