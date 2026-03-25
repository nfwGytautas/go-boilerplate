package migrators

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nfwGytautas/go-boilerplate/migrator"
)

// Postgres applies migrations to a postgres database, expects migrations to be ordered sequentially
func Postgres(ctx context.Context, dsn string, migrations []migrator.Migration) error {
	// Establish a connection
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close(ctx)

	// First check if the migration table exists and create if not
	const tableSchemaSQL = `
	CREATE TABLE IF NOT EXISTS migrations (
		id 			INT PRIMARY KEY,
		name 		VARCHAR(255) NOT NULL,
		applied_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err = db.Exec(ctx, tableSchemaSQL)
	if err != nil {
		return fmt.Errorf("failed to create migrations tables: %w", err)
	}

	// Get current migration
	const currentVersionSQL = `
	SELECT COALESCE(MAX(id), 0) FROM migrations`

	var version int
	err = db.QueryRow(ctx, currentVersionSQL).Scan(&version)
	if err == pgx.ErrNoRows {
		version = 0
	} else {
		return fmt.Errorf("failed getting current version: %w", err)
	}

	// Apply migrations
	const migrationSQL = `
	INSERT INTO migrations (id, name, applied_at) VALUES ($1, $2, $3)`

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i := 0; i < len(migrations); i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if migrations[i].Version < version {
			continue
		}

		_, err = tx.Exec(ctx, migrations[i].Content)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d, err: %w", migrations[i].Version, err)
		}

		_, err = tx.Exec(ctx, migrationSQL, migrations[i].Version, migrations[i].Name, time.Now())
		if err != nil {
			return fmt.Errorf("failed to track migration: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
