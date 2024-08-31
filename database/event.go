package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/n4ze3m/web-metrics/models"
)

func SaveEvent(ctx context.Context, db *sqlx.DB, event *models.Event) error {
	query := `
    INSERT INTO events (
        id, user_id, event_type, timestamp, url, referrer, user_agent,
        screen_width, screen_height, country, data
    ) VALUES (
        :id, :user_id, :event_type, :timestamp, :url, :referrer, :user_agent,
        :screen_width, :screen_height, :country, :data
    )`

	_, err := db.NamedExecContext(ctx, query, event)

	if err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
