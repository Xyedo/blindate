package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Healthcheck    *healthcheck
	User           *user
	BasicInfo      *basicinfo
	Location       *location
	Authentication *auth
	Tokenizer      jwtSvc
	Interest       *interest
	Online         *online
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
		ra := route.Authentication
		v1.POST("/auth", ra.postAuthHandler)
		v1.PUT("/auth", ra.putAuthHandler)
		v1.DELETE("/auth", ra.deleteAuthHandler)
	}
	{

		registerValidDObValidator()
		ru := route.User
		v1.POST("/users", ru.postUserHandler)
		auth := v1.Group("/users/:userId", validateUser(route.Tokenizer))
		{
			auth.GET("/", ru.getUserByIdHandler)
			auth.PATCH("/", ru.patchUserByIdHandler)

			ro := route.Online
			auth.POST("/online", ro.postUserOnlineHandler)
			auth.GET("/online", ro.getUserOnlineHandler)
			auth.PUT("/online", ro.putuserOnlineHandler)

			rb := route.BasicInfo
			auth.POST("/basic-info", rb.postBasicInfoHandler)
			auth.GET("/basic-info", rb.getBasicInfoHandler)
			auth.PATCH("/basic-info", rb.patchBasicInfoHandler)

			rl := route.Location
			auth.POST("/location", rl.postLocationByUserIdHandler)
			auth.GET("/location", rl.getLocationByUserIdHandler)
			auth.PATCH("/location", rl.patchLocationByUserIdHandler)

			ri := route.Interest
			auth.GET("/interests", ri.getInterestHandler)
			auth.POST("/interests/bio", ri.postInterestBioHandler)
			auth.PUT("/interests/bio", ri.putInterestBioHandler)

			interest := auth.Group("/interest/:interestId", validateInterest())
			{
				interest.POST("/hobbies", ri.postInterestHobbiesHandler)
				interest.PUT("/hobbies", ri.putInterestHobbiesHandler)
				interest.DELETE("/hobbies", ri.deleteInterestHobbiesHandler)

				interest.POST("/movie-series", ri.postInterestMovieSeriesHandler)
				interest.PUT("/movie-series", ri.putInterestMovieSeriesHandler)
				interest.DELETE("/movie-series", ri.deleteInterestMovieSeriesHandler)

				interest.POST("/travels", ri.postInterestTravelingHandler)
				interest.PUT("/travels", ri.putInterestTravelingHandler)
				interest.DELETE("/travels", ri.deleteInterestTravelingHandler)

				interest.POST("/sports", ri.postInterestSportHandler)
				interest.PUT("/sports", ri.putInterestSportHandler)
				interest.DELETE("/sports", ri.deleteInterestSportHandler)
			}
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
