package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/auth/dtos"
	"github.com/xyedo/blindate/pkg/authentication"
	"github.com/xyedo/blindate/pkg/common/app-error/httperror"
)

func NewAuth(authUC authentication.Usecase) *authH {
	return &authH{
		authUC: authUC,
	}
}

type authH struct {
	authUC authentication.Usecase
}

func (a *authH) postAuthHandler(c *gin.Context) {
	var request dtos.Login
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	accessToken, refreshToken, err := a.authUC.Login(request.Email, request.Password)
	if err != nil {
		httperror.HandleError(c, err)
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

func (a *authH) putAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	accessToken, err := a.authUC.RevalidateRefreshToken(refreshTokenCookie)
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
func (a *authH) deleteAuthHandler(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		errForbiddenResp(c, "Cookie not found in your browser, must be login")
		return
	}
	err = a.authUC.Logout(refreshTokenCookie)
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
