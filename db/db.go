package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "database/tictactoe.db")

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
		current_turn INTEGER NOT NULL,
		user_player_id INTEGER NOT NULL,
		status STRING NOT NULL,
		board BLOB NOT NULL,
		date DATETIME NOT NULL,
		FOREIGN KEY(user_owner_id) REFERENCES users(id),
		FOREIGN KEY(user_player_id) REFERENCES users(id),
		FOREIGN KEY(current_turn) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createGamesTables)
	if err != nil {
		panic("Could not create games table.")
	}

	// Enable foreign key constraints
	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		panic("Failed to enable foreign key constraints:")
	}

	query := `
INSERT INTO users (email, password, username)
SELECT ?, ?, ?
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username = ?
)`
	_, err = DB.Exec(query, "computer@local", "nologin", "computer", "computer")

	if err != nil {
		panic(err)
	}
}
