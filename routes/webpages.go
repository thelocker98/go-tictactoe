package routes

import (
	"fmt"
	"net/http"

	"example.com/tictactoe/models"
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
	games, err := models.GetGameByUserId(1)
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
		fmt.Println(tempGame, game.Board)
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
	context.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

func loadSigninPage(context *gin.Context) {
	context.HTML(http.StatusOK, "signup.html", gin.H{"title": "Signup"})
}

func test(c *gin.Context) {
	cookie, err := c.Cookie("")
	if err != nil {
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	c.String(http.StatusOK, "Cookie value: %s", cookie)
}
