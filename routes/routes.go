package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/game", loadGamePage)
	server.POST("/game", getBestMove)
}
