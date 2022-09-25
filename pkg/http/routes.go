package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/http/healthcheck"
	"github.com/xyedo/blindate/pkg/http/user"
)

type Route struct {
	Healthcheck *healthcheck.Healthcheck
	User        *user.User
}

func Routes(route Route) http.Handler {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		rh := route.Healthcheck
		v1.GET("/healthcheck", rh.HealthCheckHandler)
	}
	{
		ru := route.User
		v1.GET("/users", ru.PostUserHandler)
	}

	return r
}
