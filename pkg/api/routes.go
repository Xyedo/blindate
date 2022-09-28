package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Healthcheck *healthcheck
	User        *user
	BasicInfo   *basicinfo
}

func Routes(route Route) http.Handler {
	r := gin.New()
	r.HandleMethodNotAllowed = true

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		rh := route.Healthcheck
		v1.GET("/healthcheck", rh.healthCheckHandler)
	}

	{
		registerValidDObValidator()
		ru := route.User
		rb := route.BasicInfo
		v1.POST("/users", ru.postUserHandler)
		v1.GET("/users/:id", ru.getUserByIdHandler)
		v1.PATCH("/users/:id", ru.patchUserByIdHandler)
		v1.POST("/users/:id/basic-info", rb.postBasicInfoHandler)
		v1.GET("/users/:id/basic-info", rb.getBasicInfoHandler)
		v1.PATCH("/users/:id/basic-info", rb.patchBasicInfoHandler)
		v1.POST("/users/:id/location")
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
