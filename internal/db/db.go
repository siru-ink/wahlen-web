package db

import (
	"database/sql"
	"embed"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func InitializeDb(dbPath string) (*sql.DB, error) {
	var err error

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	err = runDbMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runDbMigrations(db *sql.DB) error {
	var err error

	goose.SetBaseFS(embeddedMigrations)

	err = goose.SetDialect("sqlite")
	if err != nil {
		return err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	return nil
}
