package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
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
		Email    string    `json:"email" binding:"required, email"`
		Password string    `json:"password" binding:"required, min=8"`
		Dob      time.Time `json:"dob" binding:"required, datetime=2000-23-08"`
	}

	c.BindJSON(&input)

	user := domain.User{
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
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": map[string]string{
			"id": user.ID,
		},
	})
}
