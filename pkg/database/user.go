package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ensureFields(user objects.User) error {
	const requiredField string = "missing field \"%s\" required when creating a new user"
	const sqlError string = "error creating new user: %w"
	if user.Username == nil {
		errMsg := fmt.Sprintf(requiredField, "username")
		return errors.New(errMsg)
	}
	if user.Password == nil {
		errMsg := fmt.Sprintf(requiredField, "password")
		return errors.New(errMsg)
	}
	return nil
}

func AddUser(db *sql.DB, newUser objects.User) error {
	err := ensureFields(newUser)
	if err != nil {
		return err
	}
	if bytes, err := bcrypt.GenerateFromPassword([]byte(*newUser.Password), 4); err != nil {
		return err
	} else {
		*newUser.Password = string(bytes)
	}
	if _, err := db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		newUser.Username, newUser.Password); err != nil {
		return err
	}
	return nil
}

// returns error and full user object for given username
func GetUser(db *sql.DB, username string) (objects.User, error) {
	var user objects.User
	err := db.QueryRow(`SELECT * FROM users WHERE username = ?`, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// queries database for all users - returns err and a slice of all users; empty if err
func GetUsers(db *sql.DB) ([]objects.User, error) {
	var users []objects.User
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		users = nil
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr objects.User
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
			users = nil
			return users, err
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		users = nil
		return users, err
	}
	return users, nil
}

// checks if user credentials are valid - returns nil on success; err otherwise
func VerifyUser(db *sql.DB, user objects.User) error {
	err := ensureFields(user)
	if err != nil {
		return err
	}
	verified, err := GetUser(db, *user.Username)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*verified.Password), []byte(*user.Password)); err != nil {
		return err
	}
	return nil
}

// delete corresponding row in users table
func DeleteUserByUsername(db *sql.DB, username string) error {
	if _, err := db.Exec(
		`DELETE FROM users WHERE username = ?;`, username); err != nil {
		return err
	}
	return nil
}

// delete corresponding row in users table
func DeleteUserById(db *sql.DB, userid int64) error {
	if _, err := db.Exec(
		`DELETE FROM users WHERE id = ?;`, userid); err != nil {
		return err
	}
	return nil
}
