package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

func (p *Postgres) FindAllShifts(ctx context.Context) ([]*models.Shift, error) {
	var shifts []*models.Shift

	query := `SELECT id, name, "in", "out", "is_default", "created_at", "updated_at" FROM shifts`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shift models.Shift

		if err := rows.Scan(
			&shift.ID,
			&shift.Name,
			&shift.In,
			&shift.Out,
			&shift.IsDefault,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		); err != nil {
			return nil, err
		}

		shifts = append(shifts, &shift)
	}

	return shifts, nil
}

func (p *Postgres) CreateShift(ctx context.Context, create models.UpsertShift) (*models.Shift, error) {
	var shift models.Shift

	query := `
		INSERT INTO shifts ("name", "in", "out", "is_default")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "name", "in", "out", "is_default", "created_at", "updated_at"`

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.
		QueryRowContext(ctx, query, create.Name, create.In, create.Out, create.IsDefault).
		Scan(&shift.ID,
			&shift.Name,
			&shift.In,
			&shift.Out,
			&shift.IsDefault,
			&shift.CreatedAt,
			&shift.UpdatedAt)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if create.IsDefault {
		if err := p.resetDefault(ctx, tx, shift.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &shift, nil
}

func (p *Postgres) FindShiftByID(ctx context.Context, id uint) (*models.Shift, error) {
	var shift models.Shift

	query := `SELECT "id", "name", "in", "out", "is_default", "created_at", "updated_at" FROM shifts WHERE "id" = $1`
	row := p.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&shift.ID,
		&shift.Name,
		&shift.In,
		&shift.Out,
		&shift.IsDefault,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.NewErrDataNotFound(id)
		}
		return nil, err
	}

	return &shift, nil
}

func (p *Postgres) UpdateShift(ctx context.Context, id uint, update models.UpsertShift) (*models.Shift, error) {
	var shift models.Shift

	query := `
		UPDATE shifts
		SET "name"=$1, "in"=$2, "out"=$3, "is_default"=$4, "updated_at"=now() WHERE id = $5
		RETURNING "id", "name", "in", "out", "is_default", "updated_at", "created_at"`

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(
		ctx,
		query,
		update.Name,
		update.In,
		update.Out,
		update.IsDefault,
		id,
	)

	if err := row.Scan(
		&shift.ID,
		&shift.Name,
		&shift.In,
		&shift.Out,
		&shift.IsDefault,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	if update.IsDefault {
		if err := p.resetDefault(ctx, tx, shift.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &shift, nil
}

func (p *Postgres) resetDefault(ctx context.Context, tx *sql.Tx, excludeID uint) error {
	query := `UPDATE shifts SET is_default=false WHERE is_default=true AND id != $1`

	_, err := tx.ExecContext(ctx, query, excludeID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteShift(ctx context.Context, id uint) error {
	result, err := p.db.ExecContext(ctx, "DELETE FROM shifts WHERE id = $1", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return db.NewErrDataNotFound(id)
	}

	return nil
}

func (p *Postgres) FindDefaultShift(ctx context.Context) (*models.Shift, error) {
	var shift models.Shift

	query := `SELECT "id", "name", "in", "out", "is_default", "created_at", "updated_at" FROM shifts WHERE "is_default"=true`
	err := p.db.
		QueryRowContext(ctx, query).
		Scan(
			&shift.ID,
			&shift.Name,
			&shift.In,
			&shift.Out,
			&shift.IsDefault,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &shift, nil
}

func (p *Postgres) FindEmployeeShift(ctx context.Context, employeeID uint, date time.Time) (*models.EmployeeShift, error) {
	var employeeShift models.EmployeeShift

	query := `SELECT id, employee_id, shift_id, date, created_at FROM employee_shifts WHERE date = $1 AND employee_id = $2`
	err := p.db.
		QueryRowContext(ctx, query, date.Format("2006-01-02"), employeeID).
		Scan(
			&employeeShift.ID,
			&employeeShift.EmployeeID,
			&employeeShift.ShiftID,
			&employeeShift.Date,
			&employeeShift.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &employeeShift, nil
}

func (p *Postgres) SaveEmployeeShift(ctx context.Context, employeeShift models.EmployeeShift) (*models.EmployeeShift, error) {
	query := `INSERT INTO employee_shifts (employee_id, shift_id, date) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := p.db.
		QueryRowContext(ctx, query, employeeShift.EmployeeID, employeeShift.ShiftID, employeeShift.Date).
		Scan(&employeeShift.ID, &employeeShift.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &employeeShift, nil
}

func (p *Postgres) DeleteEmployeeShift(ctx context.Context, unassign models.UnassignEmployeeShift) error {
	result, err := p.db.ExecContext(
		ctx,
		"DELETE FROM employee_shifts WHERE employee_id = $1 AND date = $2",
		unassign.EmployeeID,
		unassign.Date)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
