package main

import (
	"gitea.locker98.com/locker98/go-tictactoe/db"
	"gitea.locker98.com/locker98/go-tictactoe/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	server.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
