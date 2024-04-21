package models

import (
	holidayTypes "github.com/dwadp/attendance-api/internal/holiday/types"
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Holiday struct {
	ID        uint                 `json:"uint"`
	Name      string               `json:"name"`
	Type      holidayTypes.Holiday `json:"type"`
	Weekday   *time.Weekday        `json:"weekday,omitempty"`
	Date      *db.Date             `json:"date,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type UpsertHoliday struct {
	Name    string               `json:"name" validate:"required"`
	Type    holidayTypes.Holiday `json:"type" validate:"min=0,max=1"`
	Weekday holidayTypes.Weekday `json:"weekday,omitempty" validate:"required_if=Type 0,min=0,max=7"`
	Date    db.Date              `json:"date,omitempty" validate:"required_if=Type 1"`
}
