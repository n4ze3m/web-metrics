package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitPostgresDB(ctx context.Context) (*sqlx.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return db, nil
}

func createEventsTable(ctx context.Context, db *sqlx.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS website (
       id SERIAL PRIMARY KEY,
       name TEXT,
       website_url TEXT,
       website_id TEXT UNIQUE,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS events (
       id TEXT PRIMARY KEY,
       user_id TEXT,
       event_type TEXT,
       timestamp TIMESTAMP WITH TIME ZONE,
       url TEXT,
       referrer TEXT DEFAULT '',
       user_agent TEXT DEFAULT '',
       screen_width SMALLINT,
       screen_height SMALLINT,
       country CHAR(2),
       website_id TEXT,
       data JSONB,
       FOREIGN KEY (website_id) REFERENCES website(website_id)
    );
    CREATE INDEX IF NOT EXISTS idx_events_country_user_id_timestamp ON events (country, user_id, timestamp);
    `
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	if err := createEventsTable(ctx, db); err != nil {
		return err
	}
	return nil
}
