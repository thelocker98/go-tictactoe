package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "tictactoe.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		username TEXT UNIQUE NOT NULL
	);
	`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic(err)
	}

	createGamesTables := `
	CREATE TABLE IF NOT EXISTS games (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_owner_id INTEGER NOT NULL,
		user_owner_shape INTEGER NOT NULL,
		user_owner_turn_first BOOL NOT NULL,
		user_player_id INTEGER NOT NULL,
		board BLOB NOT NULL,
		date DATETIME NOT NULL,
		FOREIGN KEY(user_owner_id) REFERENCES users(id),
		FOREIGN KEY(user_player_id) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createGamesTables)
	if err != nil {
		panic("Could not create games table.")
	}

	query := "INSERT INTO users comuter@local, nologin, computer"
	_, err = DB.Exec(query)
}
