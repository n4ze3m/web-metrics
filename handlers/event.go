package handlers

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/n4ze3m/web-metrics/database"
	"github.com/n4ze3m/web-metrics/models"
	"net/http"
)

func CreateEventHandler(c echo.Context) error {
	db := c.Get("db").(*sqlx.DB)
	event := new(models.Event)
	if err := c.Bind(event); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	event.GenerateIDAndTimestamp()
	country := c.Request().Header.Get("CF-IPCountry")
	event.SetCountry(country)

	if err := event.Validate(); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := context.Background()
	if err := database.SaveEvent(ctx, db, event); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
