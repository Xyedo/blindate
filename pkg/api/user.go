package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

func NewUser(userSvc service.User) *user {
	return &user{
		userService: userSvc,
	}
}

type user struct {
	userService service.User
}

func (u *user) postUserHandler(c *gin.Context) {
	var input struct {
		FullName string    `json:"fullName" binding:"required"`
		Email    string    `json:"email" binding:"required,email"`
		Password string    `json:"password" binding:"required,min=8"`
		Dob      time.Time `json:"dob" binding:"required,validdob"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Email":    "must have an valid email",
			"Password": "must have more than 8 character",
			"Dob":      "must between today and after 1990",
		})
		if errjson != nil {
			errorServerResponse(c, err)
			return
		}
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
		if errors.Is(err, domain.ErrUniqueConstraint23505) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "email is already taken",
			})
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
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

func (u *user) getUserByIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	user, err := u.userService.GetUserById(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "id not found!",
			})
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
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
		if errors.Is(err, domain.ErrResourceNotFound) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "id not found!",
			})
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
		return
	}
	var input struct {
		FullName    *string    `json:"fullName" binding:"omitempty"`
		Email       *string    `json:"email" binding:"omitempty,email"`
		OldPassword *string    `json:"oldPassword" binding:"required_with=NewPassword,omitempty,min=8"`
		NewPassword *string    `json:"newPassword" binding:"required_with=OldPassword,omitempty,min=8"`
		Dob         *time.Time `json:"dob" binding:"omitempty,validdob"`
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Email":       "must be an valid email",
			"OldPassword": "must be more than 8 character and pairing with NewPassword",
			"NewPassword": "must be more than 8 character and pairing with OldPassword",
			"Dob":         "Must betwen today and more than 1990",
		})
		if errjson != nil {
			errorServerResponse(c, err)
			return
		}
		return
	}
	if input.NewPassword != nil && input.OldPassword != nil {
		_, err := u.userService.VerifyCredential(user.Email, *input.OldPassword)
		if err != nil {
			errorInvalidCredsResponse(c, "email or password do not match")
			return
		}
		user.Password = *input.NewPassword
	}

	if input.FullName != nil {
		user.FullName = *input.FullName
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
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "users.Id not found!")
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated",
	})

}
