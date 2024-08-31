package models

import "time"

type WebAnalytics struct {
	TotalVisitors    int           `json:"total_visitors"`
	UniqueVisitors   int           `json:"unique_visitors"`
	PageViews        int           `json:"page_views"`
	BounceRate       float64       `json:"bounce_rate"`
	AverageTimeSpent time.Duration `json:"average_time_spent"`
	TopPages         []TopPage     `json:"top_pages"`
	TopCountries     []TopCountry  `json:"top_countries"`
}

type TopPage struct {
	URL   string `json:"url"`
	Views int    `json:"views"`
}

type TopCountry struct {
	Country string `json:"country"`
	Visits  int    `json:"visits"`
}
