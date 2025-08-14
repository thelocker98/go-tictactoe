package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"gitea.locker98.com/locker98/go-tictactoe/models"
	"gitea.locker98.com/locker98/go-tictactoe/utils"
	"github.com/gin-gonic/gin"
)

func stylesCSS(context *gin.Context) {

	context.HTML(http.StatusOK, "styles.css", gin.H{})
}

func loadHomePage(context *gin.Context) {
	userId := context.GetInt64("userId")

	context.HTML(http.StatusOK, "home.html", gin.H{
		"userId": userId,
	})
}

func loadGamePage(context *gin.Context) {
	userId := context.GetInt64("userId")

	gameId, err1 := strconv.ParseInt(context.Param("id"), 10, 64)
	game, err2 := models.GetGameById(gameId)

	if err1 != nil || err2 != nil {
		context.HTML(http.StatusBadRequest, "error.html", gin.H{"ErrorMessage": "This game does not exist"})
		return
	}

	if game.UserOwnerId != userId && game.UserPlayerId != userId {
		context.HTML(http.StatusUnauthorized, "error.html", gin.H{"ErrorMessage": "You are not a player in this game"})
		return
	}

	if game.Status == "PENDING" {
		context.HTML(http.StatusOK, "error.html", gin.H{"ErrorMessage": "This Game is Pending. Please wait for your opponent to accept your game invitation"})
		return
	}
	fmt.Println(userId)
	context.HTML(http.StatusOK, "game.html", gin.H{
		"userId": userId,
	})
}

func loadLoginPage(context *gin.Context) {
	cookie, err := context.Cookie("session")
	if err != nil {
		context.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
		return
	}

	_, err = utils.VerifyToken(cookie)
	if err != nil {
		context.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
		return
	}

	context.Redirect(http.StatusFound, "/")
	context.Abort()
	return
}

func loadSigninPage(context *gin.Context) {
	cookie, err := context.Cookie("session")
	if err != nil {
		context.HTML(http.StatusOK, "signup.html", gin.H{"title": "Signup"})
		return
	}

	_, err = utils.VerifyToken(cookie)
	if err != nil {
		context.HTML(http.StatusOK, "signup.html", gin.H{"title": "Signup"})
		return
	}

	context.Redirect(http.StatusFound, "/")
	context.Abort()
	return
}

func createGamePage(context *gin.Context) {
	context.HTML(http.StatusOK, "creategame.html", gin.H{"title": "Create Game"})
}

func logout(context *gin.Context) {
	context.SetCookie("session", "", -1, "/", "", false, true)

	context.Redirect(http.StatusFound, "/login")
	context.Abort()
}
