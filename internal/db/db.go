package db

import (
	"database/sql"
	"embed"

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
