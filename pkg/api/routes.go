package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Healthcheck    *Healthcheck
	User           *User
	BasicInfo      *BasicInfo
	Location       *Location
	Authentication *Auth
	Tokenizer      jwtSvc
	Interest       *Interest
	Online         *Online
	Match          *Match
	Convo          *Conversation
	Chat           *Chat
	Webscoket      *Ws
}

func Routes(route Route) http.Handler {
	r := gin.New()
	r.HandleMethodNotAllowed = true

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.MaxMultipartMemory = 8 << 20

	registerTagName()
	registerValidDObValidator()
	registerValidEducationLevelFieldValidator()
	v1 := r.Group("/api/v1")

	rh := route.Healthcheck
	v1.GET("/healthcheck", rh.healthCheckHandler)

	ra := route.Authentication
	v1.POST("/auth", ra.postAuthHandler)
	v1.PUT("/auth", ra.putAuthHandler)
	v1.DELETE("/auth", ra.deleteAuthHandler)
	ru := route.User
	v1.POST("/users", ru.postUserHandler)
	auth := v1.Group("/", authToken(route.Tokenizer))
	user := auth.Group("/users/:userId", validateUser())
	{
		user.GET("/", ru.getUserByIdHandler)
		user.PATCH("/", ru.patchUserByIdHandler)
		user.PUT("/profile-picture", ru.putUserImageProfileHandler)

		ro := route.Online
		user.POST("/online", ro.postUserOnlineHandler)
		user.GET("/online", ro.getUserOnlineHandler)
		user.PUT("/online", ro.putuserOnlineHandler)

		rb := route.BasicInfo
		user.POST("/basic-info", rb.postBasicInfoHandler)
		user.GET("/basic-info", rb.getBasicInfoHandler)
		user.PATCH("/basic-info", rb.patchBasicInfoHandler)

		rl := route.Location
		user.POST("/location", rl.postLocationByUserIdHandler)
		user.GET("/location", rl.getLocationByUserIdHandler)
		user.PATCH("/location", rl.patchLocationByUserIdHandler)

		ri := route.Interest
		user.GET("/interests", ri.getInterestHandler)
		user.POST("/interests/bio", ri.postInterestBioHandler)
		user.PUT("/interests/bio", ri.putInterestBioHandler)

		interest := user.Group("/interest/:interestId", validateInterest())
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
	rw := route.Webscoket
	auth.GET("/ws", rw.wsEndPoint)

	rm := route.Match
	auth.GET("/new-match", rm.getNewUserToMatchHandler)
	auth.GET("/match", rm.getAllMatchRequestedHandler)
	auth.POST("/match", rm.postNewMatchHandler)

	match := auth.Group("/match/:matchId", validateMatch())
	{
		match.PUT("/request", rm.putRequestHandler)
		match.PUT("/reveal", rm.putRevealHandler)
	}

	rconv := route.Convo
	auth.POST("/conversation", rconv.postConversationHandler)
	auth.GET("/conversation", rconv.getConversationByUserId)
	conv := auth.Group("/:conversationId", validateConversation())
	{
		conv.GET("/", rconv.getConversationById)
		conv.DELETE("/", rconv.deleteConversationById)
		rchat := route.Chat
		conv.POST("/chat", rchat.postChatHandler)
		conv.POST("/chat-media", rchat.postChatMediaHandler)
		conv.GET("/chat", rchat.getMessagesHandler)
		conv.PUT("/chat/seenAt", rchat.putSeenAtHandler)

		conv.DELETE("/chat/:chatId", validateChat(), rchat.deleteMessagesByIdHandler)
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
