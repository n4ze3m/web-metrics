package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/n4ze3m/web-metrics/handlers"
)

func SetupRoutes(e *echo.Echo) {

	e.POST("/events", handlers.CreateEventHandler)

	e.GET("/stats", handlers.GetStatsHandler)
}
