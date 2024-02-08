package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
}

func New(log *slog.Logger) *Handler {

	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/image/:imageName", imageHandler)

	//TODO: implement other  endpoints

	return router
}

func imageHandler(c *gin.Context) {
	imageName := c.Param("imageName")
	c.File("./storage/image/" + imageName)
}
