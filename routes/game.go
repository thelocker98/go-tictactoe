package routes

import (
	"net/http"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/models"
	"github.com/gin-gonic/gin"
)

type webCreateGame struct {
	Opponent     int64 `json:"opponent" binding:"required"`
	Creatorfirst *bool `json:"creator_first" binding:"required"`
	Shape        int64 `json:"shape" binding:"required"`
}

type webView struct {
	GameId           int64  `json:"gameId"`
	UserOwnerName    string `json:"userOwnerName"`
	CurrentTurn      string `json:"currentTurn"`
	UserOpponentName string `json:"userOpponentName"`
	UserPlayerShape  string `json:"userPlayerShape"`
	Status           string `json:"status"`
	Win              bool   `json:"win"`
	Winner           string `json:"winner"`
	Time             string `json:"time"`
}

type homePageData struct {
	GameId           int64  `json:"gameId"`
	UserOwnerName    string `json:"userOwnerName"`
	UserOpponentName string `json:"userOpponentName"`
	UserPlayerShape  string `json:"userPlayerShape"`
	CurrentTurn      string `json:"currentTurn"`
	Win              bool   `json:"win"`
	Winner           string `json:"winner"`
	Time             string `json:"time"`
}

func createGame(context *gin.Context) {
	var webCreateGame webCreateGame
	userId := context.GetInt64("userId")

	err := context.ShouldBindJSON(&webCreateGame)
	if err != nil || userId == webCreateGame.Opponent {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	var userFirst int64
	if *webCreateGame.Creatorfirst {
		userFirst = userId
	} else {
		userFirst = webCreateGame.Opponent
	}

	newGame, err := models.NewGame(context.GetInt64("userId"), webCreateGame.Shape, userFirst, webCreateGame.Opponent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	if newGame.UserPlayerId == 1 {
		err := ai.ComputerGameAccept(newGame.GameId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"}) // Computer failed to update game
			return
		}

		// check if the it is the computers turn to go and if it is then it uses the ai to make the move
		if newGame.CurrentTurn == 1 {
			err = ai.ComputerPlayMove(newGame.GameId)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"}) // Computer failed to make a move
				return
			}
		}
	}

	context.JSON(http.StatusOK, gin.H{"gameId": newGame.GameId})
}

func userAcceptGame(data HomeWSin) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return
	}

	if game.UserPlayerId != data.UserId {
		return
	}

	if game.Status == "PENDING" {
		game.Status = "ACCEPTED"

		err := models.UpdateGame(game)
		if err != nil {
			return
		}
		return
	}
}

func userRejectGame(data HomeWSin) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return
	}

	if game.UserPlayerId != data.UserId {
		return
	}
	if game.Status == "PENDING" {
		game.Status = "DENYED"

		err := models.UpdateGame(game)

		if err != nil {
			return
		}
		return
	}
}

func userDeleteGame(data HomeWSin) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return
	}

	if game.UserOwnerId != data.UserId {
		return
	}
	if game.Status == "DENYED" || game.Status == "PENDING" {

		err := models.DeleteGame(game)

		if err != nil {
			return
		}
		return
	}
}

func loadHomepageData(userId int64) []webView {
	// Setup variables for active and past games
	var gameList []webView
	var tempGame webView

	// Get UserId
	games, err := models.GetGameByUserId(userId)
	if err != nil {
		return nil
	}

	for _, game := range games {
		win, winner := game.Board.CheckWin()

		tempGame.GameId = game.GameId

		if game.UserOwnerId == userId {
			tempGame.UserOwnerName = "You"
			player, err := models.GetUserById(game.UserPlayerId)
			if err != nil {
				continue
			}
			tempGame.UserOpponentName = player.UserName
		} else {
			tempGame.UserOpponentName = "You"
			owner, err := models.GetUserById(game.UserPlayerId)
			if err != nil {
				continue
			}
			tempGame.UserOwnerName = owner.UserName
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

		if win || tempGame.Status == "DENYED" || tempGame.Status == "DENYED_OP" {
			userShape := game.UserOwnerShape
			if userId != game.UserOwnerId {
				userShape = userShape * -1
			}
			if !win {

			} else if winner == userShape {
				tempGame.Winner = "You Won! 😃"
			} else if winner != 0 {
				tempGame.Winner = "You Lost! 😥"
			} else {
				tempGame.Winner = "You Tied! 🤝"
			}

			tempGame.Win = true
		} else {
			tempGame.Win = false
		}
		gameList = append(gameList, tempGame)
	}
	return gameList
}
