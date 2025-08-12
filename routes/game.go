package routes

import (
	"errors"
	"net/http"

	"gitea.locker98.com/locker98/go-tictactoe/ai"
	"gitea.locker98.com/locker98/go-tictactoe/models"
	"github.com/gin-gonic/gin"
)

type webCreateGame struct {
	Opponent     int64 `json:"opponent" binding:"required"`
	Creatorfirst *bool `json:"creator_first" binding:"required"`
	Shape        int64 `json:"shape" binding:"required"`
}

type webView struct {
	GameId       int64  `json:"gameId"`
	OwnerId      int64  `json:"ownerId"`
	OwnerName    string `json:"ownerName"`
	OpponentId   int64  `json:"opponentId"`
	OpponentName string `json:"opponentName"`
	CurrentTurn  int64  `json:"currentTurn"`
	OwnerShape   int64  `json:"ownerShape"`
	Win          bool   `json:"win"`
	Winner       int64  `json:"winner"`
	Status       string `json:"status"`
	Time         string `json:"time"`
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

func userAcceptGame(data gameCRUDData) (*models.Game, error) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return nil, errors.New("Faild to find the game")
	}

	if game.UserPlayerId != data.UserId {
		return nil, errors.New("User is not the player in this game")
	}

	if game.Status == "PENDING" {
		game.Status = "ACCEPTED"

		err := models.UpdateGame(game)
		if err != nil {
			return nil, errors.New("Failed to Update Database")
		}
		return game, nil
	}
	return nil, errors.New("Game in the Wrong State")
}

func userRejectGame(data gameCRUDData) (*models.Game, error) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return nil, errors.New("Faild to find the game")
	}

	if game.UserPlayerId != data.UserId {
		return nil, errors.New("User is not the player in this game")
	}
	if game.Status == "PENDING" {
		game.Status = "DENYED"

		err := models.UpdateGame(game)

		if err != nil {
			return nil, errors.New("Failed to Update Database")
		}
		return game, nil
	}
	return nil, errors.New("Game in the Wrong State")
}

func userDeleteGame(data gameCRUDData) (*models.Game, error) {
	game, err := models.GetGameById(data.GameId)

	if err != nil {
		return nil, errors.New("Faild to find the game")
	}

	if game.UserOwnerId != data.UserId {
		return nil, errors.New("User does not own this game")
	}
	if game.Status == "DENYED" || game.Status == "PENDING" {

		err := models.DeleteGame(game)

		if err != nil {
			return nil, errors.New("Delete Failed do to database error")
		}
		// add the deleted Flag so that the webpage knows not to deleted the Data
		game.Status = "DELETED"
		return game, nil
	}
	return nil, errors.New("Game in the Wrong State")
}

func loadHomepageData(userId int64) []webView {
	// Setup variables for active and past games
	var gameList []webView

	// Get UserId
	games, err := models.GetGameByUserId(userId)
	if err != nil {
		return nil
	}

	for _, game := range games {
		tempGame, err := proccessHomePageGame(game)

		if err != nil {
			continue
		}

		gameList = append(gameList, tempGame)
	}
	return gameList
}

func proccessHomePageGame(game models.Game) (webView, error) {
	win, winner := game.Board.CheckWin()

	owner, err := models.GetUserById(game.UserOwnerId)
	if err != nil {
		return webView{}, errors.New("failed to get username from DB")
	}

	opponent, err := models.GetUserById(game.UserPlayerId)
	if err != nil {
		return webView{}, errors.New("failed to get username from DB")
	}

	time := game.DateTime.Format("Jan 2 3:04 PM")

	return webView{
		GameId:       game.GameId,
		OwnerId:      game.UserOwnerId,
		OwnerName:    owner.UserName,
		OpponentId:   game.UserPlayerId,
		OpponentName: opponent.UserName,
		CurrentTurn:  game.CurrentTurn,
		OwnerShape:   game.UserOwnerShape,
		Win:          win,
		Winner:       winner,
		Status:       game.Status,
		Time:         time,
	}, nil
}
