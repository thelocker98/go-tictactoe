package routes

import (
	"example.com/tictactoe/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/login", loadLoginPage)
	server.POST("/login", login)
	server.GET("/signup", loadSigninPage)
	server.POST("/signup", signup)
	server.GET("/", loadHomePage)
	server.GET("/game", loadGamePage)
	server.GET("/test", test)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/game", getBestMove)
	authenticated.POST("/newgame", getBestMove)
	authenticated.POST("/resumegame", getBestMove)
}
