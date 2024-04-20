package models

import (
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Holiday struct {
	ID        uint          `json:"uint"`
	Name      string        `json:"name"`
	Type      int           `json:"type"`
	Weekday   *time.Weekday `json:"weekday,omitempty"`
	Date      *db.Date      `json:"date,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
