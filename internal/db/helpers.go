package db

import (
	"context"
	"database/sql"
)

func CheckUserExistence(db *sql.DB, handlerContext context.Context, username string) (exists bool) {
	row := db.QueryRowContext(handlerContext, "SELECT * FROM users WHERE username = ?", username)

	err := row.Scan() // Since this is just checking for existence no values are actually needed

	if err != nil {
		return false
	}
	return true
}

func RegisterUser(db *sql.DB, handlerContext context.Context, username, password string) (err error) {
	_, err = db.ExecContext(handlerContext, "INSERT INTO users (name, password, access) VALUES (?, ?, 'counter')", username, password)
	return err
}
