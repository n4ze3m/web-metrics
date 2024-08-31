package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID           string    `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"`
	EventType    string    `db:"event_type" json:"event_type"`
	Timestamp    time.Time `db:"timestamp" json:"timestamp"`
	URL          string    `db:"url" json:"url"`
	Referrer     string    `db:"referrer" json:"referrer"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	ScreenWidth  uint16    `db:"screen_width" json:"screen_width"`
	ScreenHeight uint16    `db:"screen_height" json:"screen_height"`
	Country      string    `db:"country" json:"country"`
	Data         string    `db:"data" json:"data,omitempty"`
	WebsiteID    string    `db:"website_id" json:"website_id"`
}

func (e *Event) GenerateIDAndTimestamp() {
	e.ID = uuid.New().String()
	e.Timestamp = time.Now().UTC()
}

func (e *Event) Validate() error {
	if e.UserID == "" {
		return fmt.Errorf("UserID is required")
	}
	if e.EventType == "" {
		return fmt.Errorf("EventType is required")
	}
	if e.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if e.ScreenWidth == 0 || e.ScreenHeight == 0 {
		return fmt.Errorf("screen dimensions are required")
	}
	if len(e.Country) != 2 {
		return fmt.Errorf("country must be a 2-letter code")
	}

	if e.WebsiteID == "" {
		return fmt.Errorf("website_id is required")
	}

	// if data is empty add a default value of "{}"

	if e.Data == "" {
		e.Data = "{}"
	}

	return nil
}

func (e *Event) SetCountry(country string) {
	if len(country) == 2 {
		e.Country = country
	} else {
		e.Country = "XX"
	}
}
