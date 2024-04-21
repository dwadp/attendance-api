package models

import (
	"github.com/dwadp/attendance-api/internal/attendance/types"
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Attendance struct {
	ID             uint                `json:"id"`
	EmployeeID     uint                `json:"employee_id"`
	ShiftID        db.NullableInt64    `json:"shift_id"`
	ShiftName      db.NullableString   `json:"shift_name"`
	ShiftIn        time.Time           `json:"shift_in"`
	ShiftOut       time.Time           `json:"shift_out"`
	ClockIn        db.NullableDateTime `json:"clock_in"`
	ClockOut       db.NullableDateTime `json:"clock_out"`
	ClockInStatus  types.Status        `json:"clock_in_status"`
	ClockOutStatus types.Status        `json:"clock_out_status"`
	Date           db.Date             `json:"date"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type AttendanceRequest struct {
	EmployeeID uint       `json:"employee_id" validate:"required"`
	Type       types.Type `json:"-"`
}
