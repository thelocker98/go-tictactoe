package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSdataIn struct {
	GameId int64 `json:"gameId"`
	Move   int64 `json:"move"`
}

var data []byte

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]int64) // Track active clients

func getBoardLayoutWS(w http.ResponseWriter, r *http.Request, userId int64) {
	// Setup and upgrade websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}

	// ensure connection is closed when done
	defer conn.Close()

	// for message recived from browser
	for {
		// read message
		_, msg, err := conn.ReadMessage()
		if err != nil {
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

		// Print username of client
		fmt.Println("User ID", userId)
		// add device to client list
		clients[conn] = GameIDData.GameId

		// Check if it is a valid move or if it is a request for games state.
		// -1 means it is a request for game state
		if GameIDData.Move == -1 {
			// get current game state from database
			webGameData, err1 := getBoardLayout(userId, GameIDData.GameId)

			var err2 error
			data, err2 = json.Marshal(webGameData)
			// check for errors in both commands
			if err1 != nil && err2 != nil {
				fmt.Println("Marshal error:", err)
				return
			}

			// print out data that will be sent
			fmt.Println(webGameData)
		} else {
			// Print out move
			fmt.Println("Move", GameIDData.Move)
			// Play move
			webGameData, err1 := playMove(userId, GameIDData)

			// check to see if it was a valid move
			var err2 error
			data, err2 = json.Marshal(webGameData)
			if err1 != nil && err2 != nil {
				fmt.Println("Marshal error:", err)
				return
			}
		}

		// Broadcast message to all clients
		fmt.Println(clients)
		for client, val := range clients {
			if val == GameIDData.GameId {
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
