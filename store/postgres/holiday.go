package postgres

import (
	"context"
	"github.com/dwadp/attendance-api/internal"
	holidayTypes "github.com/dwadp/attendance-api/internal/holiday/types"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

func (p *Postgres) FindHolidaysInDate(ctx context.Context, date time.Time) ([]*models.Holiday, error) {
	var holidays []*models.Holiday

	query := `SELECT id, name, type, weekday, date, created_at, updated_at FROM holidays WHERE date = $1 OR type = $2`
	rows, err := p.db.QueryContext(ctx, query, date.Format("2006-01-02"), holidayTypes.Weekend)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var holiday models.Holiday

		if err := rows.Scan(
			&holiday.ID,
			&holiday.Name,
			&holiday.Type,
			&holiday.Weekday,
			&holiday.Date,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
		); err != nil {
			return nil, err
		}

		holidays = append(holidays, &holiday)
	}

	return holidays, nil
}

func (p *Postgres) FindAllHoliday(ctx context.Context) ([]*models.Holiday, error) {
	var holidays []*models.Holiday

	rows, err := p.db.QueryContext(ctx, "SELECT id, name, type, weekday, date, created_at, updated_at FROM holidays ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var holiday models.Holiday

		if err := rows.Scan(
			&holiday.ID,
			&holiday.Name,
			&holiday.Type,
			&holiday.Weekday,
			&holiday.Date,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
		); err != nil {
			return nil, err
		}

		holidays = append(holidays, &holiday)
	}

	return holidays, nil
}

func (p *Postgres) FindHolidayByDateOrWeekday(ctx context.Context, date time.Time, weekday holidayTypes.Weekday) (*models.Holiday, error) {
	var holiday models.Holiday

	query := `SELECT id, name, type, weekday, date, created_at, updated_at FROM holidays WHERE date = $1`
	args := []any{date.Format("2006-01-02")}

	if weekday != holidayTypes.None {
		query += " OR weekday = $2"
		args = append(args, weekday)
	}

	err := p.db.
		QueryRowContext(ctx, query, args...).
		Scan(
			&holiday.ID,
			&holiday.Name,
			&holiday.Type,
			&holiday.Weekday,
			&holiday.Date,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &holiday, nil
}

func (p *Postgres) FindHolidayByID(ctx context.Context, id uint) (*models.Holiday, error) {
	var holiday models.Holiday

	row := p.db.QueryRowContext(ctx, "SELECT id, name, type, weekday, date, created_at, updated_at FROM holidays WHERE id = $1", id)

	if err := row.Scan(
		&holiday.ID,
		&holiday.Name,
		&holiday.Type,
		&holiday.Weekday,
		&holiday.Date,
		&holiday.CreatedAt,
		&holiday.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &holiday, nil
}

func (p *Postgres) CreateHoliday(ctx context.Context, create models.UpsertHoliday) (*models.Holiday, error) {
	var holiday models.Holiday

	row := p.db.QueryRowContext(
		ctx,
		"INSERT INTO holidays (name, type, weekday, date) VALUES ($1, $2, $3, $4) RETURNING id, name, type, weekday, date, updated_at, created_at",
		create.Name,
		create.Type,
		create.Weekday,
		create.Date,
	)

	if err := row.Scan(
		&holiday.ID,
		&holiday.Name,
		&holiday.Type,
		&holiday.Weekday,
		&holiday.Date,
		&holiday.CreatedAt,
		&holiday.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &holiday, nil
}

func (p *Postgres) UpdateHoliday(ctx context.Context, id uint, update models.UpsertHoliday) (*models.Holiday, error) {
	var holiday models.Holiday

	if update.Type == holidayTypes.NationalHoliday && update.Weekday != holidayTypes.None {
		return nil, internal.ErrNationalHolidayShouldNotHaveWeekday
	}

	if update.Type == holidayTypes.Weekend && update.Date.Valid {
		return nil, internal.ErrWeekendShouldNotHaveDate
	}

	query := `
		UPDATE holidays
		SET name=$1, type=$2, weekday=$3, date=$4, updated_at=now()
		WHERE id = $5
		RETURNING id, name, type, weekday, date, updated_at, created_at
	`
	row := p.db.QueryRowContext(
		ctx,
		query,
		update.Name,
		update.Type,
		update.Weekday,
		update.Date,
		id,
	)

	if err := row.Scan(
		&holiday.ID,
		&holiday.Name,
		&holiday.Type,
		&holiday.Weekday,
		&holiday.Date,
		&holiday.CreatedAt,
		&holiday.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &holiday, nil
}

func (p *Postgres) DeleteHoliday(ctx context.Context, id uint) error {
	result, err := p.db.ExecContext(ctx, "DELETE FROM holidays WHERE id = $1", id)
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
