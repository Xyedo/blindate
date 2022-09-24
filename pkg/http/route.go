package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/http/healthcheck"
)

func Routes() http.Handler {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/healthcheck", healthcheck.HealthCheckHandler)
	}

	return r
}
