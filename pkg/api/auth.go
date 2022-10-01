package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func NewAuth(authService service.AuthService, userService service.User, token service.Jwt) *auth {
	return &auth{
		authService: authService,
		userService: userService,
		tokenizer:   token,
	}
}

type auth struct {
	authService service.AuthService
	userService service.User
	tokenizer   service.Jwt
}

func (a *auth) postAuthHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Email":    "required and must be valid email",
			"Password": "required and must be over 8 character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	id, err := a.userService.VerifyCredential(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrNotMatchCredential) {
			errorInvalidCredsResponse(c, "email or password do not match")
			return
		}
		errorServerResponse(c, err)
		return
	}
	accessToken, err := a.tokenizer.GenerateAccessToken(id)
	if err != nil {
		errorInvalidCredsResponse(c, "fail to create token, please try again!")
		return
	}
	refreshToken, err := a.tokenizer.GenerateRefreshToken(id)
	if err != nil {
		errorInvalidCredsResponse(c, "fail to create token, please try again!")
		return
	}
	err = a.authService.AddRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrDuplicateToken) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "token is already taken, please try again",
			})
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.SetCookie("refreshToken", refreshToken, 2592000, "/api/v1", "localhost", true, true)

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": accessToken,
		},
	})
}

func (a *auth) putAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errCookieNotFound(c)
		return
	}
	err = a.authService.VerifyRefreshToken(refreshTokenCookie)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrNotMatchCredential) {
			errorInvalidCredsResponse(c, "invalid credentials")
			return
		}
		errorServerResponse(c, err)
		return
	}
	id, err := a.tokenizer.ValidateRefreshToken(refreshTokenCookie)
	if err != nil {
		if errors.Is(err, service.ErrTokenNotValid) {
			errorInvalidCredsResponse(c, "invalid credentials")
			return
		}
	}
	accessToken, err := a.tokenizer.GenerateAccessToken(id)
	if err != nil {
		errorInvalidCredsResponse(c, "fail to create token, please try again!")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": accessToken,
		},
	})

}
func (a *auth) deleteAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errCookieNotFound(c)
		return
	}
	err = a.authService.VerifyRefreshToken(refreshTokenCookie)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrNotMatchCredential) {
			errorInvalidCredsResponse(c, "invalid credentials")
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = a.authService.DeleteRefreshToken(refreshTokenCookie)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrNotMatchCredential) {
			errorInvalidCredsResponse(c, "invalid credentials")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.SetCookie("refreshToken", "", -1, "/api/v1", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "log out success",
	})
}
