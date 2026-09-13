package db

import (
	"context"
	"database/sql"
)

func CheckUserExistence(db *sql.DB, handlerContext context.Context, username string) (exists bool) {
	// row, err := db.QueryContext(handlerContext, "SELECT * FROM users WHERE username = ?", username)
	return false
}

func RegisterUser(db *sql.DB, handlerContext context.Context, username, password string) {}
