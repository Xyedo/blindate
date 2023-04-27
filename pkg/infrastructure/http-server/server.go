package httpserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/swagger"
	"github.com/xyedo/blindate/pkg/infrastructure"
)

type httpServer struct {
	Config infrastructure.Config
	Gin    *gin.Engine
}

func New(config infrastructure.Config) httpServer {
	return httpServer{
		Config: config,
		Gin:    gin.New(),
	}
}
func (h *httpServer) handler() {
	v1 := h.Gin.Group("/api/v1")

	v1.Use(cors.Default())

}
func (h *httpServer) Listen() {
	if h.Config.Env == "development" {
		h.Gin.Use(gin.Logger())
	}

	h.Gin.HandleMethodNotAllowed = true
	h.Gin.Use(gin.Recovery())

	h.Gin.GET("/swagger/*", swagger.HandlerDefault)
	h.Gin.MaxMultipartMemory = 8 << 20

	h.handler()
}
