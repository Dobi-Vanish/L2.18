package model

import "time"

// Event struct holds events.
type Event struct {
	ID     int       `json:"id,omitempty"`
	UserID int       `json:"user_id,omitempty"`
	Date   time.Time `json:"date"`
	Text   string    `json:"text"`
}
