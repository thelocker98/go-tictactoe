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
	authenticated.GET("/accept/:id", acceptGame)
	authenticated.GET("/reject/:id", rejectGame)
	authenticated.GET("/delete/:id", deleteGame)

	authenticated.GET("/creategame", createGamePage)
	authenticated.POST("/creategame", createGame)
	authenticated.GET("/users", getAllUsers)

	authenticated.GET("/game/:id", loadGamePage)
	authenticated.GET("/gamelayout/:id", getBoardLayout)
	authenticated.POST("/play", playMove)

	authenticated.GET("/logout", logout)
}
