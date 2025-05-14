package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/game", getBestMove)
	server.GET("/", loadGamePage)

}
