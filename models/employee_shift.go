package models

import (
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type EmployeeShift struct {
	ID         uint      `json:"id"`
	Date       db.Date   `json:"date"`
	EmployeeID uint      `json:"employee_id"`
	ShiftID    uint      `json:"shift_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AssignEmployeeShift struct {
	Date       db.Date `json:"date"`
	EmployeeID uint    `json:"employee_id"`
	ShiftID    uint    `json:"shift_id"`
}

type UnassignEmployeeShift struct {
	Date       db.Date `json:"date" validate:"required"`
	EmployeeID uint    `json:"employee_id" validate:"required"`
}
