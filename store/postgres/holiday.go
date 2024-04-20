package postgres

import (
	"context"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"time"
)

func (p *Postgres) FindAllHolidays(ctx context.Context, date time.Time) ([]*models.Holiday, error) {
	var holidays []*models.Holiday

	query := `SELECT id, name, type, weekday, date, created_at, updated_at FROM holidays WHERE date = $1 OR type = $2`
	rows, err := p.db.QueryContext(ctx, query, date.Format("2006-01-02"), holiday.Weekend)
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
