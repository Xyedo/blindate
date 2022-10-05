package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/service"
)

type Route struct {
	Healthcheck    *healthcheck
	User           *user
	BasicInfo      *basicinfo
	Location       *location
	Authentication *auth
	Tokenizer      service.Jwt
	Interest       *interest
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

			rb := route.BasicInfo
			auth.POST("/basic-info", rb.postBasicInfoHandler)
			auth.GET("/basic-info", rb.getBasicInfoHandler)
			auth.PATCH("/basic-info", rb.patchBasicInfoHandler)

			rl := route.Location
			auth.POST("/location", rl.postLocationByUserIdHandler)
			auth.GET("/location", rl.getLocationByUserIdHandler)
			auth.PATCH("/location", rl.patchLocationByUserIdHandler)

			ri := route.Interest
			auth.GET("/interest", ri.GetInterestHandler)
			auth.POST("/interest/bio", ri.PostInterestBioHandler)
			auth.PUT("/interest/bio", ri.PutInterestBioHandler)

			//TODO: CREATE HANDLER AND TEST INTEREST ID
			interest := auth.Group("/interest/:interestId", validateInterest())
			{
				interest.POST("/hobbies", ri.PostInterestHobbiesHandler)
				interest.PUT("/hobbies", ri.PutInterestHobbiesHandler)
				interest.DELETE("/hobbies", ri.DeleteInterestHobbiesHandler)

				interest.POST("/movie-series", ri.PostInterestMovieSeriesHandler)
				interest.PUT("/movie-series", ri.PutInterestMovieSeriesHandler)
				interest.DELETE("/movie-series", ri.DeleteInterestMovieSeriesHandler)

				interest.POST("/traveling", ri.PostInterestTravelingHandler)
				interest.PUT("/traveling", ri.PutInterestTravelingHandler)
				interest.DELETE("/traveling", ri.DeleteInterestTravelingHandler)

				interest.POST("/sport", ri.PostInterestSportHandler)
				interest.PUT("/sport", ri.PutInterestSportHandler)
				interest.DELETE("/sport", ri.DeleteInterestSportHandler)
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
