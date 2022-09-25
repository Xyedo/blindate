package user

import (
	"errors"
	"log"
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
		FullName string    `json:"fullName" binding:"required"`
		Email    string    `json:"email" binding:"required,email"`
		Password string    `json:"password" binding:"required,min=8"`
		Dob      time.Time `json:"dob" binding:"required"`
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err.Error())
	}

	user := domain.User{
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password,
		Dob:      input.Dob,
	}
	err = u.userService.CreateUser(&user)
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
		"status": "success",
		"data": map[string]string{
			"id": user.ID,
		},
	})
}
