package main

import (
	"gitea.locker98.com/locker98/go-tictactoe/db"
	"gitea.locker98.com/locker98/go-tictactoe/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	//gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.LoadHTMLGlob("templates/html/*.html")
	server.Static("/css", "templates/css")

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
