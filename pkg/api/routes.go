package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/tokenizer"
)

type Route struct {
	Healthcheck      *healthcheck
	User             *user
	BasicInfo        *basicinfo
	Location         *location
	Auththentication *auth
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
		ra := route.Auththentication
		v1.POST("/auth", ra.postAuthHandler)
		v1.PUT("/auth", ra.putAuthHandler)
		v1.DELETE("/auth", ra.deleteAuthHandler)
	}
	{

		registerValidDObValidator()
		ru := route.User
		v1.POST("/users", ru.postUserHandler)
		auth := v1.Group("/users/:id", validateUser(tokenizer.Jwt{}))
		{
			auth.GET("/", ru.getUserByIdHandler)
			auth.PATCH("/", ru.patchUserByIdHandler)

			rb := route.BasicInfo
			auth.POST("/basic-info", rb.postBasicInfoHandler)
			auth.GET("/basic-info", rb.getBasicInfoHandler)
			auth.PATCH("basic-info", rb.patchBasicInfoHandler)

			registerValidLatValidator()
			registerValidLngValidator()
			rl := route.Location
			auth.POST("location", rl.postLocationByUserIdHandler)
			auth.GET("location", rl.getLocationByUserIdHandler)
			auth.PATCH("location", rl.patchLocationByUserIdHandler)
		}

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
