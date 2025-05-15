package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loadGamePage(context *gin.Context) {
	context.HTML(http.StatusOK, "game.html", gin.H{"title": "Tic Tac Toe"})
}

func loadLoginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

func loadSigninPage(context *gin.Context) {
	context.HTML(http.StatusOK, "signup.html", gin.H{"title": "Signup"})
}
