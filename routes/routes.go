package routes

import (
	"gitea.locker98.com/locker98/go-tictactoe/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/login", loadLoginPage)
	server.POST("/login", login)
	server.POST("/signup", signup)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.GET("/", loadHomePage)

	authenticated.GET("/creategame", createGamePage)
	authenticated.POST("/creategame", createGame)
	authenticated.GET("/users", getAllUsers)

	authenticated.GET("/game/:id", loadGamePage)
	//authenticated.GET("/gamelayout/:id", getBoardLayout)
	authenticated.GET("/gamelayout/ws", func(c *gin.Context) {
		userId := c.GetInt64("userId")
		gameBoardWS(c.Writer, c.Request, userId)
	})

	authenticated.GET("/ws", func(c *gin.Context) {
		userId := c.GetInt64("userId")
		homePageWS(c.Writer, c.Request, userId)
	})
	//authenticated.POST("/play", playMove)

	authenticated.GET("/logout", logout)
}
