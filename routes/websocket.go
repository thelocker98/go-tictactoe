package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Websocket game data
type WSdataIn struct {
	GameId int64 `json:"gameId"`
	Move   int64 `json:"move"`
}

// Websocket data
var data []byte

// Websocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Websocket client list
type ClientData struct {
	GameId int64 `json:"gameId"`
	UserId int64 `json:"userId"`
}

var clients = make(map[*websocket.Conn]ClientData) // Track active clients

func getBoardLayoutWS(w http.ResponseWriter, r *http.Request, userId int64) {
	// Setup and upgrade websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}

	// ensure connection is closed when done
	defer func() {
		delete(clients, conn)
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
			delete(clients, conn)
			fmt.Println("Read error:", err)
			break
		}

		// Parse message in to the correct strucks
		var GameIDData WSdataIn
		err = json.Unmarshal(msg, &GameIDData)
		if err != nil {
			conn.Close()
			fmt.Println("error")
		}

		// add device to client list
		clients[conn] = ClientData{GameId: GameIDData.GameId, UserId: userId}

		// Check if it is a valid move or if it is a request for games state.
		// -1 means it is a request for game state
		var webGameData webGame

		if GameIDData.Move == -1 {
			// get current game state from database
			webGameData, err = getBoardLayout(userId, GameIDData.GameId)
			if err != nil {

			}
		} else {
			// Print out move
			fmt.Println("Move", GameIDData.Move)
			// Play move
			webGameData, err = playMove(userId, GameIDData)
		}
		// Check for database errors
		if err != nil {
			fmt.Println("db error:", err)
			return
		}

		// convert struct to json and check for errors
		data, err = json.Marshal(webGameData)
		if err != nil {
			fmt.Println("Marshal error:", err)
			return
		}

		// Broadcast message to all clients
		fmt.Println(clients)
		for client, val := range clients {
			if val.GameId == GameIDData.GameId {
				if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println("broadcast error:", err)
					client.Close()
					delete(clients, client)
				} else {
					fmt.Println(val)
				}
			}
		}
	}
}
