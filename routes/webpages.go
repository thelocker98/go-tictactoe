package routes

import (
	"fmt"
	"net/http"

	"example.com/tictactoe/models"
	"example.com/tictactoe/utils"
	"github.com/gin-gonic/gin"
)

type webView struct {
	GameId          int64
	UserOwnerName   string
	UserOwnerTurn   string
	UserPlayerName  string
	UserPlayerShape string
	Winner          string
}

func loadHomePage(context *gin.Context) {
	userId := context.GetInt64("userId")

	games, err := models.GetGameByUserId(userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad User"})
	}

	var activeGames []webView
	var pastGames []webView

	for _, game := range games {
		win, winner := game.Board.CheckWin()

		owner, err1 := models.GetUserById(game.UserOwnerId)
		player, err2 := models.GetUserById(game.UserPlayerId)
		if err1 != nil || err2 != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Bad User"})
		}

		var tempGame webView
		tempGame.GameId = game.GameId
		tempGame.UserOwnerName = owner.UserName
		tempGame.UserPlayerName = player.UserName

		if win {
			if winner == game.UserOwnerShape {
				tempGame.Winner = "You Won! 😃"
			} else if winner != 0 {
				tempGame.Winner = "You Lost! 😥"
			} else {
				tempGame.Winner = "You Tied! 🤝"
			}

			pastGames = append(pastGames, tempGame)
		} else {

			activeGames = append(activeGames, tempGame)
		}
	}

	fmt.Println("done")

	context.HTML(http.StatusOK, "home.html", gin.H{
		"activegames": activeGames,
		"pastgames":   pastGames,
	})
}

func loadGamePage(context *gin.Context) {
	context.HTML(http.StatusOK, "game.html", gin.H{"title": "Tic Tac Toe"})
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
