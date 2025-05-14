package webserver

import (
	"example.com/tictactoe/routes"

	"github.com/gin-gonic/gin"
)

func StartWebServer() {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
