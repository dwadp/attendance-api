package models

import (
	"encoding/json"
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

func (a *Attendance) GetClockInStatus() types.Status {
	if a.ID == 0 && !a.ClockIn.Valid && a.ShiftIn.IsZero() {
		return types.Alpha
	}

	if !a.ClockIn.Valid {
		return types.NoClockIn
	}

	return a.ClockInStatus
}

func (a *Attendance) GetClockOutStatus() types.Status {
	if a.ID == 0 && !a.ClockOut.Valid && a.ShiftOut.IsZero() {
		return types.Alpha
	}

	if !a.ClockOut.Valid {
		return types.NoClockOut
	}

	return a.ClockOutStatus
}

func (a *Attendance) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}

	value := struct {
		ID             *uint        `json:"id"`
		ShiftName      *string      `json:"shift_name"`
		ShiftIn        *string      `json:"shift_in"`
		ShiftOut       *string      `json:"shift_out"`
		ClockIn        *string      `json:"clock_in"`
		ClockOut       *string      `json:"clock_out"`
		ClockInStatus  types.Status `json:"clock_in_status"`
		ClockOutStatus types.Status `json:"clock_out_status"`
		Date           *string      `json:"date"`
		CreatedAt      time.Time    `json:"created_at"`
		UpdatedAt      time.Time    `json:"updated_at"`
	}{
		ClockInStatus:  a.GetClockInStatus(),
		ClockOutStatus: a.GetClockOutStatus(),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}

	if a.ID != 0 {
		value.ID = &a.ID
	}

	if !a.ShiftIn.IsZero() {
		shiftIn := a.ShiftIn.Format("2006-01-02 15:04:05")
		value.ShiftIn = &shiftIn
	}

	if !a.ShiftOut.IsZero() {
		shiftOut := a.ShiftOut.Format("2006-01-02 15:04:05")
		value.ShiftOut = &shiftOut
	}

	if a.ShiftName.Valid {
		shiftName := a.ShiftName.String
		value.ShiftName = &shiftName
	}

	if a.ClockIn.Valid {
		clockIn := a.ClockIn.Time.Format("2006-01-02 15:04:05")
		value.ClockIn = &clockIn
	}

	if a.ClockOut.Valid {
		clockOut := a.ClockOut.Time.Format("2006-01-02 15:04:05")
		value.ClockOut = &clockOut
	}

	if a.Date.Valid {
		date := a.Date.T.Format("2006-01-02")
		value.Date = &date
	}

	return json.Marshal(value)
}

type AttendanceRequest struct {
	EmployeeID uint       `json:"employee_id" validate:"required"`
	Type       types.Type `json:"-"`
}
