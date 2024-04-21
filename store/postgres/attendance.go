package postgres

import (
	"context"
	"fmt"
	"github.com/dwadp/attendance-api/models"
	"github.com/golang-module/carbon/v2"
	"time"
)

func (p *Postgres) SaveAttendance(ctx context.Context, attendance models.Attendance) (*models.Attendance, error) {
	query := `
		INSERT INTO
			attendances (employee_id, shift_id, shift_name, shift_in, shift_out, clock_in, clock_out, clock_in_status, clock_out_status, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	err := p.db.
		QueryRowContext(ctx, query,
			attendance.EmployeeID,
			attendance.ShiftID,
			attendance.ShiftName,
			attendance.ShiftIn,
			attendance.ShiftOut,
			attendance.ClockIn,
			attendance.ClockOut,
			attendance.ClockInStatus,
			attendance.ClockOutStatus,
			attendance.Date,
		).
		Scan(
			&attendance.ID,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (p *Postgres) FindAttendanceByEmployeeID(ctx context.Context, employeeID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance

	query := `
		SELECT
			id,
			employee_id,
			shift_id,
			shift_name,
			shift_in,
			shift_out,
			clock_in,
			clock_out,
			clock_in_status,
			clock_out_status,
			date,
			created_at,
			updated_at
		FROM attendances
		WHERE employee_id = $1 AND date = $2`

	err := p.db.
		QueryRowContext(ctx, query,
			employeeID,
			date.Format("2006-01-02"),
		).
		Scan(
			&attendance.ID,
			&attendance.EmployeeID,
			&attendance.ShiftID,
			&attendance.ShiftName,
			&attendance.ShiftIn,
			&attendance.ShiftOut,
			&attendance.ClockIn,
			&attendance.ClockOut,
			&attendance.ClockInStatus,
			&attendance.ClockOutStatus,
			&attendance.Date,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (p *Postgres) UpdateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error) {
	if attendance == nil {
		return nil, fmt.Errorf("attendance is required")
	}

	query := `
		UPDATE attendances
		SET
			clock_out = $1,
			clock_out_status = $2,
			updated_at = now()	
		WHERE id = $3
		RETURNING clock_out, clock_out_status, updated_at`

	err := p.db.
		QueryRowContext(ctx, query,
			attendance.ClockOut,
			attendance.ClockOutStatus,
			attendance.ID,
		).
		Scan(
			&attendance.ClockOut,
			&attendance.ClockOutStatus,
			&attendance.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func (p *Postgres) FindAllAttendances(ctx context.Context, employeeID uint) ([]*models.Attendance, error) {
	var attendances []*models.Attendance

	query := `
		SELECT
			id,
			employee_id,
			shift_id,
			shift_name,
			shift_in,
			shift_out,
			clock_in,
			clock_out, 
			clock_in_status,
			clock_out_status,
			date,
			created_at,
			updated_at
		FROM
			attendances
		WHERE employee_id = $1 AND (date BETWEEN $2 AND $3)
		ORDER BY date`

	rows, err := p.db.QueryContext(ctx, query,
		employeeID,
		carbon.Now().StartOfMonth().ToDateString(),
		carbon.Now().EndOfMonth().ToDateString(),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var attendance models.Attendance

		if err := rows.Scan(
			&attendance.ID,
			&attendance.EmployeeID,
			&attendance.ShiftID,
			&attendance.ShiftName,
			&attendance.ShiftIn,
			&attendance.ShiftOut,
			&attendance.ClockIn,
			&attendance.ClockOut,
			&attendance.ClockInStatus,
			&attendance.ClockOutStatus,
			&attendance.Date,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		); err != nil {
			return nil, err
		}

		attendances = append(attendances, &attendance)
	}

	return attendances, nil
}
