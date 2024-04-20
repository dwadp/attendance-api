package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store/db"
)

func (p *Postgres) FindAllShifts(ctx context.Context) ([]*models.Shift, error) {
	var shifts []*models.Shift

	query := `SELECT id, name, "in", "out", "created_at", "updated_at" FROM shifts`
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
		INSERT INTO shifts ("name", "in", "out", "created_at", "updated_at")
		VALUES ($1, $2, $3, now(), now())
		RETURNING "id", "name", "in", "out", "created_at", "updated_at"`

	err := p.db.
		QueryRowContext(ctx, query, create.Name, create.In, create.Out).
		Scan(&shift.ID,
			&shift.Name,
			&shift.In,
			&shift.Out,
			&shift.CreatedAt,
			&shift.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &shift, nil
}

func (p *Postgres) FindShiftByID(ctx context.Context, id uint) (*models.Shift, error) {
	var shift models.Shift

	query := `SELECT "id", "name", "in", "out", "created_at", "updated_at" FROM shifts WHERE "id" = $1`
	row := p.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&shift.ID,
		&shift.Name,
		&shift.In,
		&shift.Out,
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
		SET "name"=$1, "in"=$2, "out"=$3, "updated_at"=now() WHERE id = $4
		RETURNING "id", "name", "in", "out", "updated_at", "created_at"`
	row := p.db.QueryRowContext(
		ctx,
		query,
		update.Name,
		update.In,
		update.Out,
		id,
	)

	if err := row.Scan(
		&shift.ID,
		&shift.Name,
		&shift.In,
		&shift.Out,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &shift, nil
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
