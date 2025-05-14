package main

import (
	"example.com/tictactoe/db"
	"example.com/tictactoe/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	server.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
