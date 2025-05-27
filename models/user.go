package models

import (
	"errors"

	"example.com/tictactoe/db"
	"example.com/tictactoe/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
	UserName string
}

func (u *User) Save() error {
	if u.Email == "" || u.Password == "" || u.UserName == "" {
		return errors.New("missing required field")
	}
	query := "INSERT INTO users (email, password, username) VALUES (?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword, u.UserName)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return errors.New("invalid credentials")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("invalid credentials")
	}

	return nil
}

func GetUserById(userid int64) (User, error) {
	query := "SELECT id, email, username FROM users WHERE id = ?"
	row := db.DB.QueryRow(query, userid)

	var tempUser User
	err := row.Scan(&tempUser.ID, &tempUser.Email, &tempUser.UserName)

	if err != nil {
		return User{}, errors.New("invalid userid")
	}

	return tempUser, nil
}
