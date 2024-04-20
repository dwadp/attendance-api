package postgres

import (
	"context"
	"github.com/dwadp/attendance-api/models"
	"time"
)

func (p *Postgres) FindDayOff(ctx context.Context, employeeID uint, date time.Time) (*models.DayOff, error) {
	var dayOff models.DayOff

	query := `SELECT id, employee_id, description, date, created_at FROM employee_day_offs WHERE employee_id = $1 AND date = $2`
	err := p.db.
		QueryRowContext(ctx, query, employeeID, date).
		Scan(
			&dayOff.ID,
			&dayOff.EmployeeID,
			&dayOff.Description,
			&dayOff.Date,
			&dayOff.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &dayOff, nil
}

func (p *Postgres) SaveDayOff(ctx context.Context, dayOff models.DayOff) (*models.DayOff, error) {
	query := `INSERT INTO employee_day_offs (employee_id, description, date)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := p.db.QueryRowContext(ctx, query, dayOff.EmployeeID, dayOff.Description, dayOff.Date).Scan(&dayOff.ID, &dayOff.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &dayOff, nil
}
