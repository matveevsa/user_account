package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddNamedMigrationContext("20260615184754_initdb.go", upInitdb, downInitdb)
}

func upInitdb(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id BIGSERIAL PRIMARY KEY,
		    login TEXT NOT NULL UNIQUE,
		    email TEXT NOT NULL UNIQUE,
		    phone TEXT,
		    first_name TEXT,
		    second_name TEXT,
		    middle_name TEXT,
		    age INT,
		    created_at TIMESTAMP NOT NULL DEAFULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEAFULT NOW(),
		)
	`)
	return err
}

func downInitdb(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE IF EXISTS users`)
	return err
}
