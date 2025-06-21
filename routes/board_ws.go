package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameID struct {
	GameId int64 `json:"gameId"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getBoardLayoutWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close() // ensure connection is closed when done

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		var GameIDData GameID
		err = json.Unmarshal(msg, &GameIDData)
		if err != nil {
			conn.Close()
			fmt.Println("error")
		}
		fmt.Println("Parsed data:", GameIDData.GameId)

		webGameData, err1 := getBoardLayout(3, GameIDData.GameId)

		data, err2 := json.Marshal(webGameData)
		if err1 != nil && err2 != nil {
			fmt.Println("Marshal error:", err)
			return
		}

		fmt.Println(webGameData)

		conn.WriteMessage(websocket.TextMessage, data)
	}
}
