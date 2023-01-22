package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHealthCheck() *Healthcheck {
	return &Healthcheck{}
}

type Healthcheck struct{}

func (h *Healthcheck) healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "availible",
		"system_info": map[string]string{
			//todo make it use flag
			"environtment": "DEVELOPMENT",
			"version":      "0.0.1-alpha",
		},
	})
}
