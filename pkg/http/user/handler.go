package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

func New(userSvc domain.UserService) *User {
	return &User{
		userService: userSvc,
	}
}

type User struct {
	userService domain.UserService
}

func (u *User) PostUserHandler(c *gin.Context) {
	var input struct {
		FullName string    `json:"fullName" binding:"required"`
		Email    string    `json:"email" binding:"required,email"`
		Password string    `json:"password" binding:"required,min=8"`
		Dob      time.Time `json:"dob" binding:"required,validdob"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": "fail",
				"message": err1.Error(),
			})
			return
		}
		errMap := util.ReadValidationErr(err)
		if errMap != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"success": "fail",
				"message": "please refer to the documentation",
				"error":   errMap,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": "fail",
			"message": "unknown failed, fill this bug please!",
		})
		return
	}

	user := domain.User{
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password,
		Dob:      input.Dob,
	}
	err := u.userService.CreateUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "email is already taken",
			})
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "server error!",
			})
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
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
