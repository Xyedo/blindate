package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authSvc interface {
	Login(email, password string) (string, string, error)
	RevalidateRefreshToken(refreshToken string) (string, error)
	Logout(refreshToken string) error
}
type jwtSvc interface {
	GenerateAccessToken(id string) (string, error)
	GenerateRefreshToken(id string) (string, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateAccessToken(token string) (string, error)
}

func NewAuth(authService authSvc) *Auth {
	return &Auth{
		authService: authService,
	}
}

type Auth struct {
	authService authSvc
}

func (a *Auth) postAuthHandler(c *gin.Context) {
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
	accessToken, refreshToken, err := a.authService.Login(input.Email, input.Password)
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

func (a *Auth) putAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	accessToken, err := a.authService.RevalidateRefreshToken(refreshTokenCookie)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": accessToken,
		},
	})

}
func (a *Auth) deleteAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	err = a.authService.Logout(refreshTokenCookie)
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
