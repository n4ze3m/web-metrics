package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/n4ze3m/web-metrics/models"
	"time"
)

func GetWebAnalytics(websiteID string, ctx context.Context, db *sqlx.DB, startDate, endDate time.Time) (*models.WebAnalytics, error) {
	analytics := &models.WebAnalytics{}

	// If no dates are provided, default to the last 7 days
	if startDate.IsZero() && endDate.IsZero() {
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -7)
	}

	// Query for total visitors
	err := db.GetContext(ctx, &analytics.TotalVisitors, `
			SELECT COALESCE(COUNT(DISTINCT user_id), 0) 
			FROM events 
			WHERE website_id = $1 AND timestamp BETWEEN $2 AND $3`,
		websiteID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}

	// Query for unique visitors
	err = db.GetContext(ctx, &analytics.UniqueVisitors, `
			SELECT COALESCE(COUNT(DISTINCT user_id), 0) 
			FROM events 
			WHERE website_id = $1 AND event_type = 'pageview' AND timestamp BETWEEN $2 AND $3`,
		websiteID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}

	// Query for page views
	err = db.GetContext(ctx, &analytics.PageViews, `
			SELECT COALESCE(COUNT(*), 0) 
			FROM events 
			WHERE website_id = $1 AND event_type = 'pageview' AND timestamp BETWEEN $2 AND $3`,
		websiteID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}

	// Calculate bounce rate (simplified)
	var bounceCount int
	err = db.GetContext(ctx, &bounceCount, `
			SELECT COALESCE(COUNT(DISTINCT user_id), 0) 
			FROM events 
			WHERE website_id = $1 AND event_type = 'pageview' 
			AND timestamp BETWEEN $2 AND $3 
			GROUP BY user_id 
			HAVING COUNT(*) = 1`,
		websiteID, startDate, endDate,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if analytics.UniqueVisitors > 0 {
		analytics.BounceRate = float64(bounceCount) / float64(analytics.UniqueVisitors) * 100
	} else {
		analytics.BounceRate = 0
	}

	// Query for top pages
	err = db.SelectContext(ctx, &analytics.TopPages, `
			SELECT url, COUNT(DISTINCT user_id) as views 
			FROM events 
			WHERE website_id = $1 AND event_type = 'pageview' AND timestamp BETWEEN $2 AND $3 
			GROUP BY url 
			ORDER BY views DESC 
			LIMIT 10`,
		websiteID, startDate, endDate,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Query for top countries
	err = db.SelectContext(ctx, &analytics.TopCountries, `
			SELECT country, COUNT(DISTINCT user_id) as visits 
			FROM events 
			WHERE website_id = $1 AND event_type = 'pageview' AND timestamp BETWEEN $2 AND $3 
			GROUP BY country 
			ORDER BY visits DESC 
			LIMIT 10`,
		websiteID, startDate, endDate,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return analytics, nil
}
