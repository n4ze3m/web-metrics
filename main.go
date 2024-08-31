package main

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/n4ze3m/web-metrics/routes"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/n4ze3m/web-metrics/database"
)

var db *sqlx.DB

func initDB() error {
	ctx := context.Background()
	var err error
	db, err = database.InitPostgresDB(ctx)
	if err != nil {
		return err
	}
	return database.RunMigrations(ctx, db)
}

func dbMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("db", db)
		return next(c)
	}
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}(db)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(dbMiddleware)

	routes.SetupRoutes(e)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
