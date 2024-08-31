package handlers

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/n4ze3m/web-metrics/database"
)

func GetStatsHandler(c echo.Context) error {

	db := c.Get("db").(*sqlx.DB)

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start date format"})
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end date format"})
		}
	}

	// Get web analytics data
	analytics, err := database.GetWebAnalytics(c.Request().Context(), db, startDate, endDate)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch analytics data"})
	}

	return c.JSON(http.StatusOK, analytics)
}
