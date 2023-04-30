package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/domain/authentication"
	"github.com/xyedo/blindate/pkg/infrastructure"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(config infrastructure.Config, authUC authentication.Usecase) *authH {
	return &authH{
		config: config,
		authUC: authUC,
	}
}

type authH struct {
	config infrastructure.Config
	authUC authentication.Usecase
}

func (a *authH) postAuthHandler(c *gin.Context) {
	var request postAuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.Mod().Validate()
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
	refreshTokenCookie := c.GetString(constant.KeyRefreshToken)

	accessToken, err := a.authUC.RevalidateRefreshToken(refreshTokenCookie)
	if err != nil {
		httperror.HandleError(c, err)
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
	refreshTokenCookie := c.GetString(constant.KeyRefreshToken)

	err := a.authUC.Logout(refreshTokenCookie)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	c.SetCookie("refreshToken", "", -1, "/api/v1", a.config.Host, true, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "log out success",
	})
}
