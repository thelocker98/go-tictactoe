package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Websocket game data
type gameMoveData struct {
	GameId int64 `json:"gameId"`
	Move   int64 `json:"move"`
}

// Websocket game data
type gameCRUDData struct {
	UserId int64 `json:"userId"`
	GameId int64 `json:"gameId"`
	Action int64 `json:"action"` // 0 getlist, 1 accept, 2 deny, 3 delete
}

// Websocket data
var gamedata []byte
var homedata []byte

// Websocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Track active clients
var gamePageClients = make(map[*websocket.Conn]int64)
var homePageClients = make(map[*websocket.Conn]int64)

func gameBoardWS(w http.ResponseWriter, r *http.Request, userId int64) {
	// Setup and upgrade websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}

	// ensure connection is closed when done
	defer func() {
		delete(gamePageClients, conn)
		conn.Close()
	}()

	// Pingpong to checkif defices is on
	conn.SetPongHandler(func(appData string) error {
		fmt.Println("pong received")
		return nil
	})

	// for message recived from browser
	for {
		// read message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(gamePageClients, conn)
			fmt.Println("Read error:", err)
			break
		}

		// Parse message in to the correct strucks
		var GameIDData gameMoveData
		err = json.Unmarshal(msg, &GameIDData)
		if err != nil {
			conn.Close()
			fmt.Println("error")
		}

		// add device to client list
		gamePageClients[conn] = GameIDData.GameId

		// Check if it is a valid move or if it is a request for games state.
		// -1 means it is a request for game state
		var webGameData webGame

		if GameIDData.Move == -1 {
			// get current game state from database
			webGameData, err = getBoardLayout(userId, GameIDData.GameId)
			if err != nil {

			}
		} else {
			// Play move
			webGameData, err = playMove(userId, GameIDData)
		}
		// Check for database errors
		if err != nil {
			fmt.Println("db error:", err)
			return
		}

		// convert struct to json and check for errors
		gamedata, err = json.Marshal(webGameData)
		if err != nil {
			fmt.Println("Marshal error:", err)
			return
		}

		// Broadcast message to all clients
		for client, val := range gamePageClients {
			if val == GameIDData.GameId {
				if err := client.WriteMessage(websocket.TextMessage, gamedata); err != nil {
					fmt.Println("broadcast error:", err)
					client.Close()
					delete(gamePageClients, client)
				}
			}
		}
	}
}

func homePageWS(w http.ResponseWriter, r *http.Request, userId int64) {
	// Setup and upgrade websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}

	// ensure connection is closed when done
	defer func() {
		delete(homePageClients, conn)
		conn.Close()
	}()

	// Pingpong to checkif defices is on
	conn.SetPongHandler(func(appData string) error {
		fmt.Println("pong received")
		return nil
	})

	// for message recived from browser
	for {
		// read message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(homePageClients, conn)
			fmt.Println("Read error:", err)
			break
		}

		// Parse message in to the correct structs
		var homeIDData gameCRUDData
		err = json.Unmarshal(msg, &homeIDData)
		if err != nil {
			conn.Close()
			fmt.Println("error")
			return
		}
		// Make sure that userId is not trusted from websocket but retrived from cookie
		if homeIDData.UserId != userId {
			fmt.Println("userid and cookie do not match")
			conn.Close()
			return
		}

		// add device to client list
		homePageClients[conn] = userId
		var homePageData []webView

		// Look at action and make a move
		if homeIDData.Action == 0 { // 0 getlist
			homePageData = loadHomepageData(userId)
		} else if homeIDData.Action == 1 { // 1 accept
			game, err := userAcceptGame(homeIDData)
			if err != nil {
				continue
			}
			data, err := proccessHomePageGame(*game)
			if err != nil {
				continue
			}
			homePageData = append(homePageData, data)

		} else if homeIDData.Action == 2 { // 2 deny
			game, err := userRejectGame(homeIDData)
			if err != nil {
				continue
			}
			data, err := proccessHomePageGame(*game)
			if err != nil {
				continue
			}
			homePageData = append(homePageData, data)
		} else if homeIDData.Action == 3 { // 3 delete
			game, err := userDeleteGame(homeIDData)
			if err != nil {
				continue
			}
			data, err := proccessHomePageGame(*game)
			if err != nil {
				continue
			}
			homePageData = append(homePageData, data)
		}

		// convert struct to json and check for errors
		for _, data := range homePageData {
			homedata, err = json.Marshal(data)
			if err != nil {
				fmt.Println("Marshal error:", err)
				return
			}
			// Broadcast message to all clients
			for client, val := range homePageClients {
				if val == data.OwnerId || val == data.OpponentId {
					if err := client.WriteMessage(websocket.TextMessage, homedata); err != nil {
						fmt.Println("broadcast error:", err)
						client.Close()
						delete(homePageClients, client)
					}
				}
			}
		}
	}
}
