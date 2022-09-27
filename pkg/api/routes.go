package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/xyedo/blindate/pkg/api/healthcheck"
	"github.com/xyedo/blindate/pkg/api/user"
	"github.com/xyedo/blindate/pkg/domain/validation"
)

type Route struct {
	Healthcheck *healthcheck.Healthcheck
	User        *user.User
}

func Routes(route Route) http.Handler {
	r := gin.New()
	r.HandleMethodNotAllowed = true

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		rh := route.Healthcheck
		v1.GET("/healthcheck", rh.HealthCheckHandler)
	}

	{
		ru := route.User

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation("validdob", validation.ValidDob)
			if err != nil {
				panic(err)
			}
		}

		v1.POST("/users", ru.PostUserHandler)
		v1.GET("/users/:id", ru.GetUserByIdHandler)
		v1.PATCH("/users/:id", ru.PatchUserByIdHandler)
	}

	r.NoMethod(noMethod)
	r.NoRoute(noFound)
	return r
}

func noFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "failed",
		"message": "not found",
	})
}

func noMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"status":  "failed",
		"message": "method not allowed",
	})
}
