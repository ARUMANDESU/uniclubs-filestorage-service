package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Handler struct {
}

func New(log *slog.Logger) *Handler {

	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "lolkek")
	})

	//TODO: implement other  endpoints

	return router
}
