package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/tictactoe/models"
	"example.com/tictactoe/utils"
	"github.com/gin-gonic/gin"
)

type webView struct {
	GameId           int64
	UserOwnerName    string
	CurrentTurn      string
	UserOpponentName string
	UserPlayerShape  string
	Status           string
	Winner           string
	Time             string
}

func loadHomePage(context *gin.Context) {
	userId := context.GetInt64("userId")

	games, err := models.GetGameByUserId(userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad User"})
		return
	}

	var activeGames []webView
	var pastGames []webView

	for _, game := range games {
		win, winner := game.Board.CheckWin()

		owner, err1 := models.GetUserById(game.UserOwnerId)
		player, err2 := models.GetUserById(game.UserPlayerId)
		if err1 != nil || err2 != nil {
			continue
		}

		var tempGame webView
		tempGame.GameId = game.GameId

		tempGame.UserOwnerName = owner.UserName
		if game.UserOwnerId == userId {
			tempGame.UserOwnerName = "You"
			tempGame.UserOpponentName = player.UserName
		} else {
			tempGame.UserOpponentName = owner.UserName
		}

		if game.Status == "PENDING" && game.UserOwnerId == userId {
			tempGame.Status = "PENDING_OP"
		} else if game.Status == "DENYED" && game.UserOwnerId != userId {
			tempGame.Status = "DENYED_OP"
		} else {
			tempGame.Status = game.Status
		}

		if game.CurrentTurn == userId {
			tempGame.CurrentTurn = "Your Turn"
		} else {
			tempGame.CurrentTurn = "Opponents Turn"
		}

		tempGame.Time = game.DateTime.Format("Jan 2 3:04 PM")

		if win {
			userShape := game.UserOwnerShape
			if userId != game.UserOwnerId {
				userShape = userShape * -1
			}
			if winner == userShape {
				tempGame.Winner = "You Won! 😃"
			} else if winner != 0 {
				tempGame.Winner = "You Lost! 😥"
			} else {
				tempGame.Winner = "You Tied! 🤝"
			}

			pastGames = append(pastGames, tempGame)
		} else if tempGame.Status == "DENYED" || tempGame.Status == "DENYED_OP" {
			pastGames = append(pastGames, tempGame)
		} else {
			activeGames = append(activeGames, tempGame)
		}
	}

	context.HTML(http.StatusOK, "home.html", gin.H{
		"activegames": activeGames,
		"pastgames":   pastGames,
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
	context.HTML(http.StatusOK, "gamews.html", gin.H{
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
