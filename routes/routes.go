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

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.GET("/", loadHomePage)
	authenticated.GET("/game", loadGamePage)
	authenticated.POST("/game", getBestMove)
	authenticated.POST("/newgame", getBestMove)
	authenticated.POST("/resumegame", getBestMove)
	authenticated.GET("/logout", logout)
}
