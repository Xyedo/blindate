package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authSvc interface {
	AddRefreshToken(token string) error
	VerifyRefreshToken(token string) error
	DeleteRefreshToken(token string) error
}
type jwtSvc interface {
	GenerateAccessToken(id string) (string, error)
	GenerateRefreshToken(id string) (string, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateAccessToken(token string) (string, error)
}

func NewAuth(authService authSvc, userService userSvc, token jwtSvc) auth {
	return auth{
		authService: authService,
		userService: userService,
		tokenizer:   token,
	}
}

type auth struct {
	authService authSvc
	userService userSvc
	tokenizer   jwtSvc
}

func (a auth) postAuthHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"email":    "required and must be valid email",
			"password": "required and must be over 8 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	id, err := a.userService.VerifyCredential(input.Email, input.Password)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	accessToken, err := a.tokenizer.GenerateAccessToken(id)
	if err != nil {
		errUnauthorizedResp(c, "fail to create token, please try again!")
		return
	}
	refreshToken, err := a.tokenizer.GenerateRefreshToken(id)
	if err != nil {
		errUnauthorizedResp(c, "fail to create token, please try again!")
		return
	}
	err = a.authService.AddRefreshToken(refreshToken)
	if err != nil {
		jsonHandleError(c, err)
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

func (a auth) putAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	err = a.authService.VerifyRefreshToken(refreshTokenCookie)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	id, err := a.tokenizer.ValidateRefreshToken(refreshTokenCookie)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	accessToken, err := a.tokenizer.GenerateAccessToken(id)
	if err != nil {
		errUnauthorizedResp(c, "fail to create token, please try again!")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": accessToken,
		},
	})

}
func (a auth) deleteAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	err = a.authService.VerifyRefreshToken(refreshTokenCookie)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	err = a.authService.DeleteRefreshToken(refreshTokenCookie)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.SetCookie("refreshToken", "", -1, "/api/v1", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "log out success",
	})
}
