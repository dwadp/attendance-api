package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store/db"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db,
	}
}

func (p *Postgres) FindAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	var employees []*models.Employee

	rows, err := p.db.QueryContext(ctx, "SELECT id, name, phone, created_at, updated_at FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.Employee

		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Phone,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		); err != nil {
			return nil, err
		}

		employees = append(employees, &employee)
	}

	return employees, nil
}

func (p *Postgres) CreateEmployee(ctx context.Context, req models.UpsertEmployee) (*models.Employee, error) {
	var employee models.Employee

	row := p.db.QueryRowContext(
		ctx,
		"INSERT INTO employees (name, phone, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING id, name, phone, updated_at, created_at",
		req.Name,
		req.Phone,
	)

	if err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Phone,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (p *Postgres) FindEmployeeByID(ctx context.Context, id uint) (*models.Employee, error) {
	var employee models.Employee

	row := p.db.QueryRowContext(ctx, "SELECT id, name, phone, created_at, updated_at FROM employees WHERE id = $1", id)

	if err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Phone,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.NewErrDataNotFound(id)
		}
		return nil, err
	}

	return &employee, nil
}

func (p *Postgres) UpdateEmployee(ctx context.Context, id uint, update models.UpsertEmployee) (*models.Employee, error) {
	var employee models.Employee

	row := p.db.QueryRowContext(
		ctx,
		"UPDATE employees SET name=$1, phone=$2, updated_at=now() WHERE id = $3 RETURNING id, name, phone, updated_at, created_at",
		update.Name,
		update.Phone,
		id,
	)

	if err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Phone,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (p *Postgres) DeleteEmployee(ctx context.Context, id uint) error {
	result, err := p.db.ExecContext(ctx, "DELETE FROM employees WHERE id = $1", id)
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
