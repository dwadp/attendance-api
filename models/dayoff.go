package models

import (
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type DayOff struct {
	ID          uint      `json:"id"`
	EmployeeID  uint      `json:"employee_id"`
	Description string    `json:"description"`
	Date        db.Date   `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type DayOffRequest struct {
	EmployeeID  uint    `json:"employee_id" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Date        db.Date `json:"date" validate:"required"`
}
